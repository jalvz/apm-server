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

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generatePipelines(version string) {
	pipelines, err := os.Open("ingest/pipeline/definition.json")
	if err != nil {
		panic(err)
	}
	defer pipelines.Close()

	bytes, err := ioutil.ReadAll(pipelines)
	if err != nil {
		panic(err)
	}

	var definitions = make([]map[string]interface{}, 0)
	err = json.Unmarshal(bytes, &definitions)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(pipelinesPath(version), 0755)

	for _, definition := range definitions {
		pipeline, ok := definition["body"]
		if !ok {
			continue
		}
		id, ok := definition["id"]
		if !ok {
			continue
		}

		out, err := json.MarshalIndent(pipeline, "", "  ")
		if err != nil {
			panic(err)
		}
		fName := filepath.Join(pipelinesPath(version), id.(string)+".json")
		ioutil.WriteFile(fName, out, 0644)
	}
}
