package transaction

import (

	pr "github.com/elastic/apm-server/processor"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/monitoring"

	"github.com/santhosh-tekuri/jsonschema"

	"github.com/elastic/apm-server/models"
	"github.com/elastic/beats/libbeat/common"
	"time"
)

var (
	transactionMetrics = monitoring.Default.NewRegistry("apm-server.processor.transaction")
	transformations    = monitoring.NewInt(transactionMetrics, "transformations")
	validationCount    = monitoring.NewInt(transactionMetrics, "validation.count")
	validationError    = monitoring.NewInt(transactionMetrics, "validation.errors")
)

const (
	eventName          = "processor"
	processorName      = "transaction"
	transactionDocType = "transaction"
	spanDocType        = "span"
)

var schema = pr.CreateSchema(transactionSchema, processorName)

func NewProcessor(config *pr.Config) pr.Processor {
	return &processor{schema: schema, config: config}
}

type processor struct {
	schema *jsonschema.Schema
	config *pr.Config
}

func (p *processor) Validate(raw map[string]interface{}) error {
	validationCount.Inc()
	err := pr.Validate(raw, p.schema)
	if err != nil {
		validationError.Inc()
	}
	return err
}

func (p *processor) Transform(raw interface{}) ([]beat.Event, error) {
	var pa = raw.(models.TransactionsPayload)

	var events []beat.Event
	for _, tx := range pa.Transactions {
		ev := beat.Event{
			Fields: common.MapStr{
				"processor":        processorTransEntry,
				transactionDocType: common.MapStr{
					"name": tx.Name,
					"type": tx.Type,
					"result": tx.Result,
					"duration": tx.Duration,
					"sampled": tx.Sampled,
					"span_count": common.MapStr{"dropped": common.MapStr{"total": tx.SpanCount.Dropped.Total}},
					"context": common.MapStr{
						"version": tx.Context.Version,
						"custom": tx.Context.Custom,
						"tags": tx.Context.Tags,
						"request": common.MapStr{
							"method": tx.Context.Request.Method,
							"body": tx.Context.Request.Body,
						},
						"response": common.MapStr{
							"code": tx.Context.Response.StatusCode,
						},
					},
				},
				"context":          common.MapStr{
					"service": common.MapStr{
						"name": pa.Service.Name,
						"version": pa.Service.Version,
						"agent": common.MapStr{
							"name": pa.Service.Agent.Name,
							"version": pa.Service.Agent.Version,
						},
						"language": common.MapStr{
							"name": pa.Service.Language.Name,
							"version": pa.Service.Language.Version,
						},
						"framework": common.MapStr{
							"name": pa.Service.Framework.Name,
							"version": pa.Service.Framework.Version,
						},
						"runtime": common.MapStr{
							"name": pa.Service.Runtime.Name,
							"version": pa.Service.Runtime.Version,
						},
					},
					"system": common.MapStr{
						"name": pa.System.Hostname,
						"architecture": pa.System.Architecture,
						"platform": pa.System.Platform,
					},
					"process": common.MapStr{
						"pid": pa.Process.Pid,
						"ppid": pa.Process.Ppid,
						"title": pa.Process.Title,
						"argv": pa.Process.Argv,
					},
				},
			},
			Timestamp: time.Now(),
		}
		events = append(events, ev)

		trId := common.MapStr{"id": tx.ID}
		for _, sp := range tx.Spans {
			ev := beat.Event{
				Fields: common.MapStr{
					"processor":   processorSpanEntry,
					spanDocType:   common.MapStr{
						"type": sp.Type,
						"duration": sp.Duration,
						"name": sp.Name,
						"parent": sp.Parent,
						"start": sp.Start,
						"id": sp.ID,
						"stacktrace": getStacktrace(sp.Stacktrace),
					},
					"transaction": trId,
					"context":     common.MapStr{
						"db": common.MapStr{
							"user": sp.Context.Db.User,
							"statement": sp.Context.Db.Statement,
							"type": sp.Context.Db.Type,
							"instance": sp.Context.Db.Instance,
						},
					},
				},
				Timestamp: time.Now(),
			}
			events = append(events, ev)
		}
	}
	return events, nil
}

func (p *processor) Name() string {
	return processorName
}


func getStacktrace(st models.SpanStacktrace) []common.MapStr {
	rs := make([]common.MapStr, 0)
	for _, frame := range st {
		rs = append(rs, common.MapStr{
			"context_line":frame.ContextLine,
			"colno":frame.Colno,
			"module":frame.Module,
			"library_frame":frame.LibraryFrame,
			"abs_path":frame.AbsPath,
			"function":frame.Function,
			"post_context":frame.PostContext,
			"pre_contect":frame.PreContext,
			"vars":frame.Vars,
			"lineno":frame.Lineno,
			"filename": frame.Filename,
		})
	}
	return rs
}
