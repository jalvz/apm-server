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

package config

import "github.com/elastic/beats/v7/libbeat/common"

func NewIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{
		Meta:       &Meta{Package: &Package{}},
		DataStream: &OuterDataStream{},
		Streams: []Stream{{
			DataStream: &InnerDataStream{},
		}},
	}
}

// This represents the whole configuration that comes from Elastic Agent
// it will need to be updated along with the apm-server spec in Elastic Agent, eg. to add the `fleet` key.
type IntegrationConfig struct {
	Id         string           `config:"id"`
	Name       string           `config:"name"`
	Revision   int              `config:"revision"`
	Type       string           `config:"type"`
	UseOutput  string           `config:"use_output"`
	Meta       *Meta            `config:"meta"`
	DataStream *OuterDataStream `config:"data_stream"`
	Streams    []Stream         `config:"streams"`
}

type Stream struct {
	Id         string           `config:"id"`
	DataStream *InnerDataStream `config:"data_stream"`
	// TODO this will be lifted up to IntegrationConfig when https://github.com/elastic/package-spec/issues/70 lands
	APMServer *common.Config `config:"apm-server"`
}

type OuterDataStream struct {
	// Used in index suffix
	Namespace string `config:"namespace"`
}

// there are 2 data_stream keys with different structure, so we need 2 different types for them - hence Outer / Inner
type InnerDataStream struct {
	// Until https://github.com/elastic/package-spec/issues/64 is implemented, we probably should use this value for the middle index segment
	Dataset string `config:"dataset"`
	Type    string `config:"type"`
}

type Meta struct {
	Package *Package `config:"package"`
}

type Package struct {
	Name string `config:"name"`
	// Version will be very important for maintaining backwards-compatibility
	Version string `config:"version"`
}
