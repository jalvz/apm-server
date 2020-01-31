package stream

import (
	"time"

	"github.com/elastic/apm-server/model/metadata"
	"github.com/elastic/apm-server/publish"
)

type Decoder interface {
	Decode(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error)
}

type DecoderFunc func(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error)

func (f DecoderFunc) Decode(input interface{}, requestTime time.Time, metadata metadata.Metadata) (publish.Transformable, error) {
	return f(input, requestTime, metadata)
}
