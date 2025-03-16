package ktcenter

import (
	"sync"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/aiagt/kitextool/log"

	"github.com/kitex-contrib/config-nacos/nacos"
)

type NacosConfigCenter struct {
	opts      *nacos.Options
	client    nacos.Client
	callbacks []Callback

	once sync.Once
}

func WithNacosConfigCenter(opts *nacos.Options) ConfigCenter {
	return &NacosConfigCenter{opts: opts}
}

func (c *NacosConfigCenter) Client() nacos.Client {
	if c.client == nil {
		panic("the nacos client is not initialized")
	}

	return c.client
}

func (c *NacosConfigCenter) Initialize(conf *ktconf.CenterConf) {
	c.once.Do(func() {
		if c.opts == nil {
			c.opts = new(nacos.Options)
		}

		opts := c.opts

		if conf != nil {
			utils.SetDefault(&opts.Address, conf.GetHost())
			utils.SetDefault(&opts.Port, uint64(conf.Port))
			utils.SetDefault(&opts.NamespaceID, conf.Key)
		}

		if opts.ConfigParser == nil {
			opts.ConfigParser = NewStringParser[vo.ConfigType]()
		}

		client, err := nacos.NewClient(*opts)
		if err != nil {
			panic(err)
		}

		c.client = client
	})
}

func (c *NacosConfigCenter) RegisterCallbacks(callbacks ...Callback) {
	c.callbacks = append(c.callbacks, callbacks...)
}

func (c *NacosConfigCenter) Register(dest string, conf ktconf.Conf) {
	param, err := c.Client().ServerConfigParam(&nacos.ConfigParamConfig{
		Category:          dynamicConfigName,
		ServerServiceName: dest,
	})
	if err != nil {
		log.Fatal(err)
	}

	uniqueID := nacos.GetUniqueID()

	server.RegisterShutdownHook(func() {
		if err := c.Client().DeregisterConfig(param, uniqueID); err != nil {
			log.Fatal(err)
		}
	})

	c.Client().RegisterConfigCallback(param, func(data string, parser nacos.ConfigParser) {
		err := ktconf.ParseConf(data, conf)
		if err != nil {
			log.Errorf("parse conf failed: %s", err.Error())
			return
		}

		for _, callback := range c.callbacks {
			callback(conf)
		}
	}, uniqueID)
}
