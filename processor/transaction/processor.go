package transaction

import (

	pr "github.com/elastic/apm-server/processor"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/monitoring"

	"github.com/santhosh-tekuri/jsonschema"

	"github.com/elastic/apm-server/utility"
	"time"
	"errors"
	"fmt"
	"github.com/elastic/apm-server/model"

	"github.com/olivere/elastic"

	"github.com/elastic/apm-server/client"
	"runtime"
	"encoding/json"
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

var TransactionBuffer chan payload

const BufferSize = 1000

func init() {
	TransactionBuffer = make(chan payload, BufferSize)
	bulkProcessor := client.BulkProcessor(processorName, runtime.NumCPU())
	go consume(bulkProcessor)
}

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
	pa := decode(raw.(map[string]interface{}))
	select {
	case TransactionBuffer <- pa:
		fmt.Println("PUSHED!")
		return nil, nil
	case <- time.After(time.Second):
		fmt.Println("BLOCKED!")
		return nil, errors.New("QUEUE FULL")
	}
}

func (p *processor) Name() string {
	return processorName
}

func consume(bulk *elastic.BulkProcessor) {
	for pa := range TransactionBuffer {
		var bars [][]byte
		fmt.Println("CONSUME!")
		for _, e := range pa.Events {
			eba, err := json.Marshal(e)
			if err != nil {
				panic(err)
			}
			bars = append(bars, eba)

			for _, s := range e.Spans {
				sba, err := json.Marshal(s)
				if err != nil {
					panic(err)
				}
				bars = append(bars, sba)
			}
		}
		client.SaveToES(bars, bulk)
	}
}

func decode(raw map[string]interface{}) payload {

	pa := payload{Process: &model.Process{}, System: &model.System{},
		Service: model.Service{
			Agent:     model.Agent{},
			Framework: model.Framework{},
			Language:  model.Language{},
			Runtime:   model.Runtime{}}}
	pa.Service.Name, _ = utility.DeepGet(raw, "service")["name"].(string)
	serviceVersion, _ := utility.DeepGet(raw, "service")["version"].(string)
	pa.Service.Version = &serviceVersion
	pa.Service.Agent.Name, _ = utility.DeepGet(raw, "service", "agent")["name"].(string)
	pa.Service.Agent.Version, _ = utility.DeepGet(raw, "service", "agent")["version"].(string)
	frameworkName, _ := utility.DeepGet(raw, "service", "framework")["name"].(string)
	pa.Service.Framework.Name = &frameworkName
	frameworkVersion, _ := utility.DeepGet(raw, "service", "framework")["version"].(string)
	pa.Service.Framework.Version = &frameworkVersion
	runtimeName, _ := utility.DeepGet(raw, "service", "runtime")["name"].(string)
	pa.Service.Runtime.Name = &runtimeName
	runtimeVersion, _ := utility.DeepGet(raw, "service", "runtime")["version"].(string)
	pa.Service.Runtime.Version = &runtimeVersion
	languageName, _ := utility.DeepGet(raw, "service", "language")["name"].(string)
	pa.Service.Language.Name = &languageName
	languageVersion, _ := utility.DeepGet(raw, "service", "language")["version"].(string)
	pa.Service.Language.Version = &languageVersion
	environment, _ := utility.DeepGet(raw, "service")["environment"].(string)
	pa.Service.Environment = &environment

	pa.Process.Pid, _ = utility.DeepGet(raw, "process")["pid"].(int)
	ppid, _ := utility.DeepGet(raw, "process")["ppid"].(int)
	pa.Process.Ppid = &ppid
	title, _ := utility.DeepGet(raw, "process")["title"].(string)
	pa.Process.Title = &title
	pa.Process.Argv, _ = utility.DeepGet(raw, "process")["argv"].([]string)

	hostname, _ := utility.DeepGet(raw, "system")["hostname"].(string)
	pa.System.Hostname = &hostname
	platform, _ := utility.DeepGet(raw, "system")["platform"].(string)
	pa.System.Platform = &platform
	architecture, _ := utility.DeepGet(raw, "system")["architecture"].(string)
	pa.System.Architecture = &architecture

	if txs, ok := raw["transactions"].([]interface{}); ok {
		pa.Events = make([]Event, len(txs))
		for txIdx, tx := range txs {
				tx, ok := tx.(map[string]interface{})
				if !ok {
						continue
					}
				event := Event{}
				event.Context = utility.DeepGet(tx, "context")
				event.Id, _ = tx["id"].(string)
				name, _ := tx["name"].(string)
				event.Name = &name
				event.Type, _ = tx["type"].(string)
				event.Duration, _ = tx["duration"].(float64)
				result, _ := tx["result"].(string)
				event.Result = &result
				if timestamp, ok := tx["timestamp"].(string); ok {
						event.Timestamp, _ = time.Parse(time.RFC3339, timestamp)
					}
				event.Sampled, _ = tx["sampled"].(*bool)
				//droppedTotal, _ := utility.DeepGet(tx, "span_count", "dropped")["total"].(int)
				//event.Dropped.Total = &droppedTotal
				event.Marks = utility.DeepGet(tx, "marks")
				if spans, ok := tx["spans"].([]interface{}); ok {
						event.Spans = make([]*Span, len(spans))
						span := Span{}
						for spIdx, sp := range spans {
								sp, ok := sp.(map[string]interface{})
								if !ok {
										continue
									}
								id, _ := sp["id"].(int)
								span.Id = &id
								parent, _ := sp["parent"].(int)
								span.Parent = &parent
								span.Name, _ = sp["name"].(string)
								span.Type, _  = sp["type"].(string)
								span.Start, _ = sp["start"].(float64)
								span.Duration, _ = sp["duration"].(float64)
								span.Context = utility.DeepGet(sp, "context")
								if frames, ok := sp["stacktrace"].([]interface{}); ok {
										span.Stacktrace = make(model.Stacktrace, len(frames))
										for frIdx, fr := range frames {
												fr, ok := fr.(map[string]interface{})
												if !ok {
														continue
													}
												frame := model.StacktraceFrame{}
												function, _  := fr["function"].(string)
												frame.Function = &function
												absPath, _ := fr["abs_path"].(string)
												frame.AbsPath = &absPath
												frame.Filename, _ = fr["filename"].(string)
												frame.Lineno, _ = fr["lineno"].(int)
												libraryFrame, _ := fr["library_frame"].(bool)
												frame.LibraryFrame = &libraryFrame
												frame.Vars = utility.DeepGet(fr, "vars")
												module, _ := fr["module"].(string)
												frame.Module = &module
												colno, _ := fr["colno"].(int)
												frame.Colno = &colno
												contextLine, _ := fr["context_line"].(string)
												frame.ContextLine = &contextLine
												frame.PreContext, _ = fr["pre_context"].([]string)
												frame.PostContext, _ = fr["post_context"].([]string)
												span.Stacktrace[frIdx] = &frame
											}
									} else {
										span.Stacktrace = make(model.Stacktrace, 0)
									}
								event.Spans[spIdx] = &span
							}
					} else {
						event.Spans = make([]*Span, 0)
					}
				pa.Events[txIdx] = event
			}
	} else {
		pa.Events = make([]Event, 0)
	}
	return pa
}
