package error

import (
	"github.com/santhosh-tekuri/jsonschema"

	"time"

	"github.com/pkg/errors"

	"github.com/elastic/apm-server/model"
	pr "github.com/elastic/apm-server/processor"
	"github.com/elastic/apm-server/utility"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/monitoring"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/elastic/apm-server/client"
	"encoding/json"
)

var (
	errorMetrics    = monitoring.Default.NewRegistry("apm-server.processor.error")
	validationCount = monitoring.NewInt(errorMetrics, "validation.count")
	validationError = monitoring.NewInt(errorMetrics, "validation.errors")
	transformations = monitoring.NewInt(errorMetrics, "transformations")
)

const (
	processorName = "error"
	errorDocType  = "error"
)

var schema = pr.CreateSchema(errorSchema, processorName)

var ErrorBuffer chan payload

const BufferSize = 50

func init() {
	ErrorBuffer = make(chan payload, BufferSize)
	bulkProcessor := client.BulkProcessor(processorName, 1)
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
	case ErrorBuffer <- pa:
		return nil, nil
	case <-time.After(time.Second):
		fmt.Println("BLOCKED")
		return nil, errors.New("QUEUE FULL")
	}

}

func (p *processor) Name() string {
	return processorName
}

func consume(bulk *elastic.BulkProcessor) {
	for pa := range ErrorBuffer {
		var bars [][]byte
		fmt.Println("CONSUME!")
		for _, e := range pa.Events {
			eba, err := json.Marshal(e)
			if err != nil {
				panic(err)
			}
			bars = append(bars, eba)
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

	if errs, ok := raw["errors"].([]interface{}); ok {
		pa.Events = make([]Event, len(errs))
		for errIdx, err := range errs {
			err, ok := err.(map[string]interface{})
			if !ok {
				continue
			}
			event := Event{Log: &Log{}, Exception: &Exception{}, Transaction: &Transaction{}}
			event.Context = utility.DeepGet(err, "context")
			id, _ := err["id"].(string)
			event.Id = &id
			if timestamp, ok := err["timestamp"].(string); ok {
				event.Timestamp, _ = time.Parse(time.RFC3339, timestamp)
			}
			culprit, _ := err["culprit"].(string)
			event.Culprit = &culprit

			log := utility.DeepGet(err, "log")
			event.Log.Message, _ = log["message"].(string)
			paramMessage, _ := log["param_message"].(string)
			event.Log.ParamMessage = &paramMessage
			loggerName, _ := log["logger_name"].(string)
			event.LoggerName = &loggerName
			logLevel, _ := log["level"].(string)
			event.Log.Level = &logLevel

			if frames, ok := log["stacktrace"].([]interface{}); ok {
				event.Log.Stacktrace = make(model.Stacktrace, len(frames))
				for frIdx, fr := range frames {
					fr, ok := fr.(map[string]interface{})
					if !ok {
						continue
					}
					frame := model.StacktraceFrame{}
					function, _ := fr["function"].(string)
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
					event.Log.Stacktrace[frIdx] = &frame
				}
			} else {
				event.Log.Stacktrace = make(model.Stacktrace, 0)
			}

			event.Transaction.Id, _ = utility.DeepGet(err, "transaction")["id"].(string)
			ex := utility.DeepGet(err, "exception")
			event.Exception.Message, _ = ex["message"].(string)
			exType, _ := ex["type"].(string)
			event.Exception.Type = &exType
			exModule, _ := ex["module"].(string)
			event.Exception.Module = &exModule
			event.Exception.Code = ex["code"]
			exHandled, _ := ex["handled"].(bool)
			event.Exception.Handled = &exHandled
			event.Exception.Attributes = ex["attributes"]

			if frames, ok := ex["stacktrace"].([]interface{}); ok {
				event.Exception.Stacktrace = make(model.Stacktrace, len(frames))
				for frIdx, fr := range frames {
					fr, ok := fr.(map[string]interface{})
					if !ok {
						continue
					}
					frame := model.StacktraceFrame{}
					function, _ := fr["function"].(string)
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
					event.Exception.Stacktrace[frIdx] = &frame
				}
			} else {
				event.Exception.Stacktrace = make(model.Stacktrace, 0)
			}

			pa.Events[errIdx] = event
		}
	} else {
		pa.Events = make([]Event, 0)
	}
	return pa
}
