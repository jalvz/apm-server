// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package otlp_test

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/elastic/apm-server/beater/otlp"
	"github.com/elastic/apm-server/publish"
	"github.com/elastic/apm-server/transform"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/monitoring"
)

var (
	exportTraceServiceRequestType  = proto.MessageType("opentelemetry.proto.collector.trace.v1.ExportTraceServiceRequest")
	exportTraceServiceResponseType = proto.MessageType("opentelemetry.proto.collector.trace.v1.ExportTraceServiceResponse")
)

func TestConsumeTraces(t *testing.T) {
	var events []transform.Transformable
	var reportError error
	report := func(ctx context.Context, req publish.PendingReq) error {
		events = append(events, req.Transformables...)
		return reportError
	}

	// Send a minimal trace to verify that everything is connected properly.
	//
	// We intentionally do not check the published event contents; those are
	// tested in processor/otel.
	cannedRequest := jsonExportTraceServiceRequest(`{
"resource_spans": [
  {
    "instrumentation_library_spans": [
      {
        "spans": [
	  {
	    "trace_id": "0123456789abcdef0123456789abcdef",
	    "span_id": "945254c567a5417e",
	    "name": "operation_name"
	  }
	]
      }
    ]
  }
]
}`)

	conn := newServer(t, report)
	err := conn.Invoke(
		context.Background(), "/opentelemetry.proto.collector.trace.v1.TraceService/Export",
		cannedRequest, newExportTraceServiceResponse(),
	)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	reportError = errors.New("failed to publish events")
	err = conn.Invoke(
		context.Background(), "/opentelemetry.proto.collector.trace.v1.TraceService/Export",
		cannedRequest, newExportTraceServiceResponse(),
	)
	assert.Error(t, err)
	errStatus := status.Convert(err)
	assert.Equal(t, "failed to publish events", errStatus.Message())
	assert.Len(t, events, 2)

	actual := map[string]interface{}{}
	monitoring.GetRegistry("apm-server.otlp.grpc.consumer").Do(monitoring.Full, func(key string, value interface{}) {
		actual[key] = value
	})
	assert.Equal(t, map[string]interface{}{
		"request.count":         int64(2),
		"response.count":        int64(2),
		"response.errors.count": int64(1),
		"response.valid.count":  int64(1),
	}, actual)
}

func jsonExportTraceServiceRequest(j string) interface{} {
	request := reflect.New(exportTraceServiceRequestType.Elem()).Interface()
	decoder := json.NewDecoder(strings.NewReader(j))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(request); err != nil {
		panic(err)
	}
	return request
}

func newExportTraceServiceResponse() interface{} {
	return reflect.New(exportTraceServiceResponseType.Elem()).Interface()
}

func newServer(t *testing.T, report publish.Reporter) *grpc.ClientConn {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	srv := grpc.NewServer()
	err = otlp.RegisterGRPCServices(srv, report, logp.NewLogger("otlp_test"))
	require.NoError(t, err)

	go srv.Serve(lis)
	t.Cleanup(srv.GracefulStop)
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return conn
}
