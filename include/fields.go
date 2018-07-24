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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package include

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("apm-server", "fields.yml", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
