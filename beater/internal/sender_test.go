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

package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/apm-server/server"
)

func TestOkBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "_", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	Send(w, req, server.OK(map[string]interface{}{"some": "body"}))
	rsp := w.Result()
	got := body(t, rsp)
	assert.Equal(t, "{\"some\":\"body\"}\n", string(got))
	assert.Equal(t, "text/plain; charset=utf-8", rsp.Header.Get("Content-Type"))
}

func TestOkBodyJson(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "_", nil)
	req.Header.Set("Accept", "application/json")
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	Send(w, req, server.OK(map[string]interface{}{"version": "1.0"}))
	rsp := w.Result()
	got := body(t, rsp)
	assert.Equal(t,
		`{
  "version": "1.0"
}
`, string(got))
	assert.Equal(t, "application/json", rsp.Header.Get("Content-Type"))
}

func TestAccept(t *testing.T) {
	expectedErrorJson :=
		`{
  "error": "error message"
}
`
	expectedErrorText := "{\"error\":\"error message\"}\n"

	for idx, test := range []struct{ accept, expectedError, expectedContentType string }{
		{"application/json", expectedErrorJson, "application/json"},
		{"*/*", expectedErrorJson, "application/json"},
		{"text/html", expectedErrorText, "text/plain; charset=utf-8"},
		{"", expectedErrorText, "text/plain; charset=utf-8"},
	} {
		req, err := http.NewRequest(http.MethodPost, "_", nil)
		require.NoError(t, err)
		if test.accept != "" {
			req.Header.Set("Accept", test.accept)
		} else {
			delete(req.Header, "Accept")
		}
		w := httptest.NewRecorder()
		Send(w, req, server.BadRequest(errors.New("error message")))
		rsp := w.Result()
		got := body(t, rsp)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, test.expectedError, got, fmt.Sprintf("at index %d", idx))
		assert.Equal(t, test.expectedContentType, rsp.Header.Get("Content-Type"), fmt.Sprintf("at index %d", idx))
	}
}

func TestIncCounter(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "_", nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	e := errors.New("")
	for i := 1; i <= 5; i++ {
		for _, res := range []server.Response{
			server.EmptyResponse(http.StatusAccepted),
			server.OK(nil),
			server.Forbidden(e),
			server.Unauthorized(),
			server.RequestTooLarge(e),
			server.RateLimited(),
			server.MethodNotAllowed(),
			server.BadRequest(e),
			server.FullQueue(e),
		} {
			SendCnt(w, req, res)
			assert.Equal(t, int64(i), counterMap[res.Code()].Get())
		}
	}
	assert.Equal(t, int64(45), respCounter.Get())
	assert.Equal(t, int64(10), validCounter.Get())
	assert.Equal(t, int64(35), errCounter.Get())
}

func body(t *testing.T, response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)
	return string(body)
}
