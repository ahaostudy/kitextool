package ktconf

import (
	"github.com/aiagt/kitextool/utils"
	"github.com/aiagt/kitextool/utils/lists"
)

type ServerConf struct {
	Server       Server     `yaml:"server"`
	Log          Log        `yaml:"log"`
	Registry     Registry   `yaml:"registry"`
	ConfigCenter CenterConf `yaml:"config_center"`
	DB           *DB        `yaml:"db"`
	DBs          []*DB      `yaml:"dbs"`
	Redis        *Redis     `yaml:"redis"`
	Redises      []*Redis   `yaml:"redises"`
	Rabbitmq     Rabbitmq   `yaml:"rabbitmq"`
}

func (d *ServerConf) ParseServerConf(data string) error {
	return Parse([]byte(data), d)
}

func (d *ServerConf) GetServerConf() *ServerConf {
	return d
}

func (d *ServerConf) ParseClientConf(data string) error {
	return nil
}

func (d *ServerConf) GetClientConf(_ string) *ClientConf {
	return nil
}

type Server struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Version string `yaml:"version"`
}

type Log struct {
	Level         string `yaml:"level"`
	FileName      string `yaml:"file_name"`
	MaxSize       int    `yaml:"max_size"`
	MaxBackups    int    `yaml:"max_backups"`
	MaxAge        int    `yaml:"max_age"`
	EnableFile    *bool  `yaml:"enable_file"`
	FlushInterval int    `yaml:"flush_interval"`
}

type Registry struct {
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

func (c *Registry) GetAddress() []string {
	return lists.Map(c.Address, utils.CompleteAddress)
}

type CenterConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}

func (c *CenterConf) GetHost() string {
	const localhost = "localhost"

	if c.Host == "" {
		return localhost
	}

	return c.Host
}

type DB struct {
	Name string `yaml:"name"`
	DSN  string `yaml:"dsn"`
}

type Redis struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Rabbitmq struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"vhost"`
}

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

func (l Log) LogLevel() LogLevel {
	switch l.Level {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "notice":
		return LevelNotice
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}
