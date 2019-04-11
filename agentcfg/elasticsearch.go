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

package agentcfg

type m = map[string]interface{}

type ms = []m

type Result struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Hits []ActualHit `json:"hits"`
}

type ActualHit struct {
	Id     string `json:"_id"`
	Source Source `json:"_source"`
}

type Source struct {
	Settings m `json:settings`
}

func nameEnvQuery(name, env string) m {
	return m{
		"query": m{
			"bool": m{
				"must": ms{
					m{
						"term": m{
							"service.name": name,
						},
					},
					m{
						"term": m{
							"service.environment": env,
						},
					},
				},
			},
		},
		"sort": m{
			"@timestamp": m{
				"order": "desc",
			},
		},
		"size": 1,
	}
}

func nameQuery(name string) m {
	return m{
		"query": m{
			"bool": m{
				"must": m{
					"term": m{
						"service.name": name,
					},
				},
				"must_not": m{
					"exists": m{
						"field": "service.environment",
					},
				},
			},
		},
		"sort": m{
			"@timestamp": m{
				"order": "desc",
			},
		},
		"size": 1,
	}
}
