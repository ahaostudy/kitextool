package ktregistry

import (
	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/log"
	"github.com/cloudwego/kitex/pkg/registry"
	consulapi "github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
)

func NewConsulRegistry() Registry {
	return func(conf *ktconf.Registry) []registry.Registry {
		result := make([]registry.Registry, len(conf.GetAddress()))

		for _, address := range conf.GetAddress() {
			consulConfig := consulapi.Config{
				Address:  address,
				HttpAuth: &consulapi.HttpBasicAuth{Username: conf.Username, Password: conf.Password},
			}

			r, err := consul.NewConsulRegisterWithConfig(&consulConfig)
			if err != nil {
				log.Fatalf("service registry failed: %s", err.Error())
			}

			result = append(result, r)
		}

		return result
	}
}
