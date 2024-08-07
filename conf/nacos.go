package ktconf

import (
	"sync"

	"github.com/ahaostudy/kitextool/log"

	"github.com/kitex-contrib/config-nacos/nacos"
)

type NacosConfigCenter struct {
	opts      nacos.Options
	client    nacos.Client
	callbacks []Callback

	once sync.Once
}

func NewNacosConfigCenter(opts nacos.Options) *NacosConfigCenter {
	c := &NacosConfigCenter{opts: opts}
	return c
}

func (c *NacosConfigCenter) Client() nacos.Client {
	if c.client == nil {
		panic("the nacos client is not initialized")
	}
	return c.client
}

func (c *NacosConfigCenter) Init(conf *CenterConf) {
	c.once.Do(func() {
		opts := c.opts
		if conf != nil {
			if opts.Address == "" {
				opts.Address = conf.Host
			}
			if opts.Port == 0 {
				opts.Port = uint64(conf.Port)
			}
			if opts.NamespaceID == "" {
				opts.NamespaceID = conf.Key
			}
		}
		if opts.ConfigParser == nil {
			opts.ConfigParser = DefaultNacosParser()
		}
		client, err := nacos.NewClient(opts)
		if err != nil {
			panic(err)
		}
		c.client = client
	})
}

func (c *NacosConfigCenter) RegisterCallbacks(callbacks ...Callback) {
	c.callbacks = callbacks
}

func (c *NacosConfigCenter) Register(dest string, conf Conf) {
	param, err := c.Client().ServerConfigParam(&nacos.ConfigParamConfig{
		Category:          dynamicConfigName,
		ServerServiceName: dest,
	})
	if err != nil {
		panic(err)
	}
	c.Client().RegisterConfigCallback(param, func(data string, parser nacos.ConfigParser) {
		err := ParseConf([]byte(data), conf)
		if err != nil {
			log.Errorf("parse conf failed: %s", err.Error())
			return
		}
		for _, callback := range c.callbacks {
			callback(conf)
		}
	}, nacos.GetUniqueID())
}
