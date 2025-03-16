package ktcenter

import (
	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/log"
)

var dynamicConfigName = "config"

type Callback func(conf ktconf.Conf)

type ConfigCenter interface {
	Initialize(conf *ktconf.CenterConf)
	RegisterCallbacks(callbacks ...Callback)
	Register(dest string, conf ktconf.Conf)
}

func ApplyDynamicConfig(center ConfigCenter, centerConf *ktconf.CenterConf, dest string, conf ktconf.Conf) {
	logger := func(c ktconf.Conf) {
		log.Infof("config changed: %+v", c)
	}

	center.Initialize(centerConf)
	center.RegisterCallbacks(logger)
	center.Register(dest, conf)
}
