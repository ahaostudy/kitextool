package ktresolver

import (
	"net"
	"strconv"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/log"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosResolver(conf *ktconf.Resolver) []discovery.Resolver {
	result := make([]discovery.Resolver, len(conf.GetAddress()))

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
			log.Fatalf("service resolver failed: %s", err.Error())
		}

		result = append(result, resolver.NewNacosResolver(cli))
	}

	return result
}
