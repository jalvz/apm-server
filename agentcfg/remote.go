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

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/kibana"
	"github.com/elastic/beats/libbeat/outputs/elasticsearch"
)

type RemoteSource interface {
	Fetch(map[string]string) (map[string]interface{}, string, error)
}

type ValidationError error

type noopRemote struct{}

func NoopRemote() RemoteSource {
	return noopRemote{}
}

func (_ noopRemote) Fetch(_ map[string]string) (map[string]interface{}, string, error) {
	return nil, "", nil
}

type kibanaRemote struct {
	*kibana.Client
}

func (k kibanaRemote) Fetch(_ map[string]string) (map[string]interface{}, string, error) {
	return nil, "", errors.New("not implemented")
}

func KibanaRemote(cfg *common.Config) (RemoteSource, error) {
	kibana, error := kibana.NewKibanaClient(cfg)
	return &kibanaRemote{kibana}, error
}

type elasticsearchRemote struct {
	*elasticsearch.Client
}

func ElasticSearchRemote(cfg *common.Config) (RemoteSource, error) {
	es, err := elasticsearch.NewConnectedClient(cfg)
	return &elasticsearchRemote{es}, err
}

func (es elasticsearchRemote) Fetch(filters map[string]string) (map[string]interface{}, string, error) {
	name, _ := filters[ServiceName]
	if name == "" {
		return nil, "", ValidationError(errors.New(ServiceName + " is required"))
	}
	env, _ := filters[ServiceEnv]
	var q map[string]interface{}
	if env == "" {
		q = nameQuery(name)
	} else {
		q = nameEnvQuery(name, env)
	}
	bytes, err := es.request(q)

	if err == nil {
		res := &Result{}
		err = json.Unmarshal(bytes, res)
		if len(res.Hits.Hits) == 1 {
			hit := res.Hits.Hits[0]
			return hit.Source.Settings, hit.Id, nil
		}
	}
	return nil, "", err
}

func (es elasticsearchRemote) request(q map[string]interface{}) ([]byte, error) {
	_, ret, err := es.Request(http.MethodPost, "/cm/_search", "", nil, q)
	return ret, err
}
