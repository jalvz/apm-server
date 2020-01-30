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

package converter

import (
	"regexp"

	"github.com/elastic/apm-server/model"
	"github.com/elastic/apm-server/model/error"
	"github.com/elastic/apm-server/model/metricset"
	"github.com/elastic/apm-server/model/onboarding"
	"github.com/elastic/apm-server/model/profile"
	"github.com/elastic/apm-server/model/span"
	"github.com/elastic/apm-server/model/transaction"
	"github.com/elastic/apm-server/sourcemap"
	"github.com/elastic/beats/libbeat/beat"
)

// EventConverter knows how to convert from APM Event to beat Event types
type EventConverter struct {
	libraryPattern      *regexp.Regexp
	excludeFromGrouping *regexp.Regexp
	sourcemapStore      *sourcemap.Store
}

func NewConverter(libraryPattern, excludeFromGrouping string, sourcemapStore *sourcemap.Store) EventConverter {
	return EventConverter{
		libraryPattern:      regexp.MustCompile(libraryPattern),
		excludeFromGrouping: regexp.MustCompile(excludeFromGrouping),
		sourcemapStore:      sourcemapStore,
	}
}

func (converter EventConverter) ToBeatEvents(event model.Transformable) []beat.Event {
	switch e := event.(type) {
	case transaction.Event:
		return e.Transform()
	case span.Event:
		return e.Transform(converter.libraryPattern, converter.excludeFromGrouping, converter.sourcemapStore)
	case error.Event:
		return e.Transform(converter.libraryPattern, converter.excludeFromGrouping, converter.sourcemapStore)
	case metricset.Metricset:
		return e.Transform()
	case profile.PprofProfile:
		return e.Transform()
	case onboarding.OnboardingDoc:
		return e.Transform()
	default:
		return []beat.Event{}
	}
}
