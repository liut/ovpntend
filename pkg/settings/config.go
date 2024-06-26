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
	ManageNames []string `envconfig:"manage_names" default:"first"`
	StatusDir   string   `envconfig:"status_dir"`

	CookieName   string `envconfig:"Cookie_Name" default:"staff"`
	CookiePath   string `envconfig:"Cookie_Path" default:"/"`
	CookieDomain string `envconfig:"Cookie_Domain"`
	CookieMaxAge int    `envconfig:"Cookie_MaxAge"`

	EasyRSABIN        string `envconfig:"EASYRSA_BIN" default:"/usr/local/bin"`
	EasyRSAPKI        string `envconfig:"EASYRSA_PKI" default:"/etc/openvpn/pki"`
	EasyRSACaPassword string `envconfig:"EASYRSA_CA_PASSWORD"`

	ServerHost  string `envconfig:"SERVER_NAME"`
	ServerPort  int    `envconfig:"SERVER_PORT" default:"1194"`
	ServerProto string `envconfig:"SERVER_PROTO" default:"udp"`
	ServerPlace string `envconfig:"SERVER_Place" `

	MailEnabled        bool   `envconfig:"SMTP_ENABLED"`
	MailHost           string `envconfig:"SMTP_HOST"`
	MailPort           int    `envconfig:"SMTP_PORT" default:"465"`
	MailSenderName     string `envconfig:"SMTP_SENDER_NAME" default:"notify"`
	MailSenderEmail    string `envconfig:"SMTP_SENDER_EMAIL"`
	MailSenderPassword string `envconfig:"SMTP_SENDER_PASSWORD"`
	MailTLSEnabled     bool   `envconfig:"SMTP_TLS" default:"true"`

	ValidMailDomains []string `envconfig:"VALID_MAIL_DOMAINS"`

	RedisURI string `envconfig:"REDIS_URI" default:"redis://localhost:6379/1"`
}

func init() {
	_ = envconfig.Process(Name, Current)
	Current.Name = Name
}

// Usage 显示配置内容
func Usage() {
	_ = envconfig.Usage(Name, Current)
}
