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

package beater

import (
	"net/http"

	"github.com/elastic/apm-server/agentcfg"
)

func agentConfigHandler(remote agentcfg.RemoteSource, secretToken string) http.Handler {

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		params := r.URL.Query()
		cfg, etag, err := remote.Fetch(
			map[string]string{
				agentcfg.ServiceName: params.Get(agentcfg.ServiceName),
				agentcfg.ServiceEnv:  params.Get(agentcfg.ServiceEnv),
			})

		if err != nil {
			switch err.(type) {
			case agentcfg.ValidationError:
				send(w, r, toJson(err), http.StatusBadRequest)
			default:
				send(w, r, toJson(err), http.StatusInternalServerError)
			}
			return
		}

		if len(cfg) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Cache-Control", "max-age=0")
		if clientEtag := r.Header.Get("If-None-Match"); clientEtag != "" && clientEtag == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if etag != "" {
			w.Header().Set("Etag", etag)
		}

		send(w, r, cfg, http.StatusOK)
	})

	return authHandler(secretToken, logHandler(handler))
}
