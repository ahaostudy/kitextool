package ktregistry

import (
	"net"
	"strconv"

	"github.com/aiagt/kitextool/log"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/registry"
	nacosregistry "github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosRegistry() Registry {
	return func(conf *ktconf.Registry) []registry.Registry {
		result := make([]registry.Registry, len(conf.GetAddress()))

		for _, address := range conf.GetAddress() {
			host, portStr, err := net.SplitHostPort(address)
			if err != nil {
				log.Fatal(err)
			}

			if host == "" {
				host = "127.0.0.1"
			}

			port, err := strconv.ParseUint(portStr, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			sc := []constant.ServerConfig{
				*constant.NewServerConfig(host, port),
			}
			cc := constant.ClientConfig{
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				Username:            conf.Username,
				Password:            conf.Password,
			}

			cli, err := clients.NewNamingClient(
				vo.NacosClientParam{
					ClientConfig:  &cc,
					ServerConfigs: sc,
				},
			)
			if err != nil {
				log.Fatalf("service registry failed: %s", err.Error())
			}

			result = append(result, nacosregistry.NewNacosRegistry(cli))
		}

		return result
	}
}
