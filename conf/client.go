package ktconf

import (
	"github.com/aiagt/kitextool/utils"
	"github.com/aiagt/kitextool/utils/lists"
)

type MultiClientConf map[string]ClientConf

func (m *MultiClientConf) ParseClientConf(data string) error {
	if m == nil {
		*m = make(MultiClientConf)
	}

	return Parse([]byte(data), m)
}

func (m *MultiClientConf) GetClientConf(name string) *ClientConf {
	if c, ok := (*m)[name]; ok {
		return &c
	}

	return nil
}

func (m *MultiClientConf) ParseServerConf(data string) error {
	return nil
}

func (m *MultiClientConf) GetServerConf() *ServerConf {
	return nil
}

type ClientConf struct {
	Resolver Resolver `yaml:"resolver"`
}

func (c *ClientConf) ParseClientConf(data string) error {
	return Parse([]byte(data), c)
}

func (c *ClientConf) GetClientConf(_ string) *ClientConf {
	return c
}

func (c *ClientConf) ParseServerConf(data string) error {
	return nil
}

func (c *ClientConf) GetServerConf() *ServerConf {
	return nil
}

type Resolver struct {
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

func (c *Resolver) GetAddress() []string {
	return lists.Map(c.Address, utils.CompleteAddress)
}
