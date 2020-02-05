package settings

import (
	"github.com/kelseyhightower/envconfig"
)

// Current 当前配置
var (
	Current = new(Config)
	Name    = "ovpn"
)

// Config 配置结构
type Config struct {
	Name        string   `ignored:"true"`
	HTTPListen  string   `envconfig:"http_listen" default:":7580"`
	SentryDSN   string   `envconfig:"sentry_dsn"`
	ManageAddrs []string `envconfig:"manage_addrs" default:"127.0.0.1:7505"` // 127.0.0.1:7504,127.0.0.1:7505
	StatusDir   string   `envconfig:"status_dir"`
}

func init() {
	_ = envconfig.Process(Name, Current)
	Current.Name = Name
}

// Usage 显示配置内容
func Usage() {
	envconfig.Usage(Name, Current)
}
