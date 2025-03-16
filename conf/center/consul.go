package ktcenter

import (
	"fmt"
	"sync"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/utils"
	"github.com/cloudwego/kitex/server"

	"github.com/aiagt/kitextool/log"

	"github.com/kitex-contrib/config-consul/consul"
)

type ConsulConfigCenter struct {
	opts      *consul.Options
	client    consul.Client
	callbacks []Callback

	once sync.Once
}

func WithConsulConfigCenter(opts *consul.Options) ConfigCenter {
	return &ConsulConfigCenter{opts: opts}
}

func (c *ConsulConfigCenter) Client() consul.Client {
	if c.client == nil {
		panic("the consul client is not initialized")
	}

	return c.client
}

func (c *ConsulConfigCenter) Initialize(conf *ktconf.CenterConf) {
	c.once.Do(func() {
		if c.opts == nil {
			c.opts = new(consul.Options)
		}

		opts := c.opts

		if conf != nil {
			utils.SetDefault(&opts.Addr, fmt.Sprintf("%s:%d", conf.GetHost(), conf.Port))
		}

		if opts.ConfigParser == nil {
			opts.ConfigParser = NewStringParser[consul.ConfigType]()
		}

		client, err := consul.NewClient(*c.opts)
		if err != nil {
			panic(err)
		}

		c.client = client
	})
}

func (c *ConsulConfigCenter) RegisterCallbacks(callbacks ...Callback) {
	c.callbacks = callbacks
}

func (c *ConsulConfigCenter) Register(dest string, conf ktconf.Conf) {
	param, err := c.Client().ServerConfigParam(&consul.ConfigParamConfig{
		Category:          dynamicConfigName,
		ServerServiceName: dest,
	})
	if err != nil {
		log.Fatal(err)
	}

	var (
		key      = param.Prefix + "/" + param.Path
		uniqueID = consul.AllocateUniqueID()
	)

	c.Client().RegisterConfigCallback(
		key, uniqueID,
		func(data string, parser consul.ConfigParser) {
			err := ktconf.ParseConf(data, conf)
			if err != nil {
				log.Errorf("parse conf failed: %s", err.Error())
				return
			}

			for _, callback := range c.callbacks {
				callback(conf)
			}
		})

	server.RegisterShutdownHook(func() {
		c.Client().DeregisterConfig(key, uniqueID)
	})
}
