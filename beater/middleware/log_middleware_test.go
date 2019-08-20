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

package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/elastic/apm-server/beater/beatertest"
	"github.com/elastic/apm-server/beater/headers"
	"github.com/elastic/apm-server/beater/request"
	logs "github.com/elastic/apm-server/log"
)

func TestLogMiddleware(t *testing.T) {
	err := logp.DevelopmentSetup(logp.ToObserverOutput())
	require.NoError(t, err)

	testCases := []struct {
		name, message string
		level         zapcore.Level
		handler       request.Handler
		code          int
	}{
		{
			name:    "Accepted",
			message: "handled request",
			level:   zapcore.InfoLevel,
			handler: beatertest.Handler202,
			code:    http.StatusAccepted,
		},
		{
			name:    "Error",
			message: "error handling request",
			level:   zapcore.ErrorLevel,
			handler: beatertest.Handler403,
			code:    http.StatusForbidden,
		},
		{
			name:    "Panic",
			message: "error handling request",
			level:   zapcore.ErrorLevel,
			handler: RecoverPanicMiddleware()(beatertest.HandlerPanic),
			code:    http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := beatertest.DefaultContextWithResponseRecorder()
			c.Request.Header.Set(headers.UserAgent, tc.name)
			LogMiddleware()(tc.handler)(c)
			assert.Equal(t, tc.code, rec.Code)
			for i, entry := range logp.ObserverLogs().TakeAll() {
				// expect only one log entry per request
				assert.Equal(t, i, 0)
				assert.Equal(t, logs.Request, entry.LoggerName)
				assert.Equal(t, tc.level, entry.Level)
				assert.Equal(t, tc.message, entry.Message)

				ec := entry.ContextMap()
				assert.NotEmpty(t, ec["request_id"])
				assert.NotEmpty(t, ec["method"])
				assert.Equal(t, c.Request.URL.String(), ec["URL"])
				assert.NotEmpty(t, ec["remote_address"])
				assert.Equal(t, c.Request.Header.Get(headers.UserAgent), ec["user-agent"])
				// zap encoded type
				assert.Equal(t, tc.code, int(ec["response_code"].(int64)))
			}
		})
	}
}
