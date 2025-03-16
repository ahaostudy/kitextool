package ktresolver

import (
	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/log"
	"github.com/cloudwego/kitex/pkg/discovery"
	consulapi "github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
)

func NewConsulResolver(conf *ktconf.Resolver) []discovery.Resolver {
	result := make([]discovery.Resolver, len(conf.GetAddress()))

	for _, address := range conf.GetAddress() {
		consulConfig := consulapi.Config{
			Address:  address,
			HttpAuth: &consulapi.HttpBasicAuth{Username: conf.Username, Password: conf.Password},
		}

		r, err := consul.NewConsulResolverWithConfig(&consulConfig)
		if err != nil {
			log.Fatalf("service resolver failed: %s", err.Error())
		}

		result = append(result, r)
	}

	return result
}
