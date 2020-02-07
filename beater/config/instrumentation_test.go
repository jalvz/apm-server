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

import (
	"testing"

	"github.com/elastic/go-ucfg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/beats/libbeat/common"
)

func TestNonzeroHosts(t *testing.T) {
	t.Run("ZeroHost", func(t *testing.T) {
		cfg, err := NewConfig("9.9.9",
			common.MustNewConfigFrom(map[string]interface{}{
				"instrumentation.enabled": true, "instrumentation.hosts": []string{""}}), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), ucfg.ErrZeroValue.Error())
		assert.Nil(t, cfg)
	})

	t.Run("Valid", func(t *testing.T) {
		cfg, err := NewConfig("9.9.9",
			common.MustNewConfigFrom(map[string]interface{}{"instrumentation.enabled": true}), nil)
		require.NoError(t, err)
		assert.True(t, *cfg.SelfInstrumentation.Enabled)
		assert.Empty(t, cfg.SelfInstrumentation.Hosts)
	})
}
