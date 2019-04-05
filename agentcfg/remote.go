package agentcfg

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/kibana"
	"github.com/elastic/beats/libbeat/outputs/elasticsearch"
)

type RemoteSource interface {
	Fetch() map[string]map[string]interface{}
}

type noopRemote struct{}

func NoopRemote() RemoteSource {
	return noopRemote{}
}

func (_ noopRemote) Fetch() map[string]map[string]interface{} {
	return nil
}

type kibanaRemote struct {
	*kibana.Client
}

func (k kibanaRemote) Fetch() map[string]map[string]interface{} {
	return nil
}

func KibanaRemote(cfg *common.Config) (RemoteSource, error) {
	kibana, error := kibana.NewKibanaClient(cfg)
	return &kibanaRemote{kibana}, error
}

type elasticsearchRemote struct {
	*elasticsearch.Client
}

func (es elasticsearchRemote) Fetch() map[string]map[string]interface{} {
	return nil
}

func ElasticSearchRemote(cfg *common.Config) (RemoteSource, error) {
	es, err := elasticsearch.NewConnectedClient(cfg)
	return &elasticsearchRemote{es}, err
}
