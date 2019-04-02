package elasticsearch

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/go-elasticsearch"
	"net/url"
)

type Config struct {
	Hosts    []string `config:"hosts"`
	Protocol string   `config:"protocol"`
	Username string   `config:"username"`
	Password string   `config:"password"`
}

func FromBeatsConfig(beatCfg *common.Config) (*elasticsearch.Client, error) {
	esCfg := Config{}
	if err := beatCfg.Unpack(&esCfg); err != nil {
		return nil, err
	}
	addresses := make([]string, 0)
	for _, host := range esCfg.Hosts {
		url, _ := url.Parse(host)
		url.Scheme = esCfg.Protocol
		addresses = append(addresses, url.String())
	}
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: addresses,
		Username:  esCfg.Username,
		Password:  esCfg.Password,
	})
	return c, err
}
