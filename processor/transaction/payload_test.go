package transaction

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"time"

	m "github.com/elastic/apm-server/model"
	pr "github.com/elastic/apm-server/processor"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/apm-server/tests/loader"

	"encoding/json"
)

func TestPayloadTransform(t *testing.T) {
	hostname := "a.b.c"
	architecture := "darwin"
	platform := "x64"
	timestamp := time.Now()

	service := m.Service{Name: "myservice"}
	system := &m.System{
		Hostname:     &hostname,
		Architecture: &architecture,
		Platform:     &platform,
	}

	txValid := Event{Timestamp: timestamp}
	txValidEs := common.MapStr{
		"context": common.MapStr{
			"service": common.MapStr{
				"name":  "myservice",
				"agent": common.MapStr{"name": "", "version": ""},
			},
		},
		"processor": common.MapStr{
			"event": "transaction",
			"name":  "transaction",
		},
		"transaction": common.MapStr{
			"duration": common.MapStr{"us": 0},
			"id":       "",
			"type":     "",
			"sampled":  true,
		},
	}

	txValidWithSystem := common.MapStr{
		"processor": common.MapStr{
			"event": "transaction",
			"name":  "transaction",
		},
		"transaction": common.MapStr{
			"duration": common.MapStr{"us": 0},
			"id":       "",
			"type":     "",
			"sampled":  true,
		},
		"context": common.MapStr{
			"system": common.MapStr{
				"hostname":     hostname,
				"architecture": architecture,
				"platform":     platform,
			},
			"service": common.MapStr{
				"name":  "myservice",
				"agent": common.MapStr{"name": "", "version": ""},
			},
		},
	}
	txWithContext := Event{Timestamp: timestamp, Context: common.MapStr{"foo": "bar", "user": common.MapStr{"id": "55"}}}
	txWithContextEs := common.MapStr{
		"processor": common.MapStr{
			"event": "transaction",
			"name":  "transaction",
		},
		"transaction": common.MapStr{
			"duration": common.MapStr{"us": 0},
			"id":       "",
			"type":     "",
			"sampled":  true,
		},
		"context": common.MapStr{
			"foo": "bar", "user": common.MapStr{"id": "55"},
			"service": common.MapStr{
				"name":  "myservice",
				"agent": common.MapStr{"name": "", "version": ""},
			},
			"system": common.MapStr{
				"hostname":     "a.b.c",
				"architecture": "darwin",
				"platform":     "x64",
			},
		},
	}
	spans := []*Span{{}}
	txValidWithSpan := Event{Timestamp: timestamp, Spans: spans}
	spanEs := common.MapStr{
		"context": common.MapStr{
			"service": common.MapStr{
				"name":  "myservice",
				"agent": common.MapStr{"name": "", "version": ""},
			},
		},
		"processor": common.MapStr{
			"event": "span",
			"name":  "transaction",
		},
		"span": common.MapStr{
			"duration": common.MapStr{"us": 0},
			"name":     "",
			"start":    common.MapStr{"us": 0},
			"type":     "",
		},
		"transaction": common.MapStr{"id": ""},
	}

	tests := []struct {
		Payload Payload
		Output  []common.MapStr
		Msg     string
	}{
		{
			Payload: Payload{Service: service, Events: []Event{}},
			Output:  nil,
			Msg:     "Payload with empty Event Array",
		},
		{
			Payload: Payload{
				Service: service,
				Events:  []Event{txValid, txValidWithSpan},
			},
			Output: []common.MapStr{txValidEs, txValidEs, spanEs},
			Msg:    "Payload with multiple Events",
		},
		{
			Payload: Payload{
				Service: service,
				System:  system,
				Events:  []Event{txValid},
			},
			Output: []common.MapStr{txValidWithSystem},
			Msg:    "Payload with System and Event",
		},
		{
			Payload: Payload{
				Service: service,
				System:  system,
				Events:  []Event{txWithContext},
			},
			Output: []common.MapStr{txWithContextEs},
			Msg:    "Payload with Service, System and Event with context",
		},
	}

	for idx, test := range tests {
		outputEvents := test.Payload.transform(&pr.Config{})
		for j, outputEvent := range outputEvents {
			assert.Equal(t, test.Output[j], outputEvent.Fields, fmt.Sprintf("Failed at idx %v; %s", idx, test.Msg))
			assert.Equal(t, timestamp, outputEvent.Timestamp)
		}

	}
}

func TestPayload_MarshalJSON(t *testing.T) {
	byt, _ := loader.LoadValidDataAsBytes("transaction")
	var pa Payload

	json.Unmarshal(byt, &pa)
	fmt.Println(pa)
	for _, e := range pa.Events {

		byt2, err := json.Marshal(e)

		assert.Nil(t, err)

		var m map[string]interface{}

		json.Unmarshal(byt2, &m)

		fmt.Println(m)
		fmt.Println(len(byt2))
	}



}
