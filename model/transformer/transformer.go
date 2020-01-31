package transformer

import (
	"regexp"
	"time"

	"github.com/elastic/apm-server/model/metadata"
	"github.com/elastic/apm-server/publish"
	"github.com/elastic/apm-server/sourcemap"
)

// Transformer encapsulates configuration for decoding events into
// model objects, which can later be transformed into beat.Events.
type Transformer struct {
	Experimental        bool
	LibraryPattern      *regexp.Regexp
	ExcludeFromGrouping *regexp.Regexp
	SourcemapStore      *sourcemap.Store
}

func (t *Transformer) DecodeTransaction(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error) {
	//return transaction.Decoder(t.Experimental)
	panic("TODO")
}

func (t *Transformer) DecodeSpan(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error) {
	//return span.Decoder(t.Experimental)
	panic("TODO")
}

func (t *Transformer) DecodeError(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error) {
	//return apmerror.Decoder(t.Experimental)
	panic("TODO")
}

func (t *Transformer) DecodeMetricset(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error) {
	//return metricset.Decode()
	panic("TODO")
}
