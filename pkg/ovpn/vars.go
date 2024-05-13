package ovpn

import (
	"fmt"
	"os"
	"strings"

	"github.com/liut/ovpntend/pkg/settings"
)

// nolint
var (
	easyrasPKI string
	caPassword string

	cmdEnv []string

	smtpHost   string
	smtpPort   int
	smtpUser   string
	smtpPass   string
	senderName string

	serverHost  string
	serverPort  int
	serverProto string
	serverPlace string

	ssuffix string
)

func init() {
	easyrasPKI = settings.Current.EasyRSAPKI
	caPassword = settings.Current.EasyRSACaPassword
	cmdEnv = append(os.Environ(),
		"EASYRSA="+settings.Current.EasyRSABIN,
		"EASYRSA_PKI="+settings.Current.EasyRSAPKI,
		"OVPN_CN="+settings.Current.ServerHost,
		"OVPN_DEFROUTE=0",
		fmt.Sprintf("OVPN_PORT=%d", settings.Current.ServerPort),
		"OVPN_PROTO="+settings.Current.ServerProto,
	)

	smtpHost = settings.Current.MailHost
	smtpPort = settings.Current.MailPort
	smtpUser = settings.Current.MailSenderEmail
	smtpPass = settings.Current.MailSenderPassword
	senderName = settings.Current.MailSenderName

	serverHost = settings.Current.ServerHost
	serverPort = settings.Current.ServerPort
	serverProto = settings.Current.ServerProto
	serverPlace = settings.Current.ServerPlace

	if serverPlace != "" {
		ssuffix = "-" + serverPlace
	}
}

func IsValidName(name string) bool {
	if 0 == len(name) {
		return false
	}
	return IsValidEmail(name)
}

func IsValidEmail(email string) bool {
	if len(settings.Current.ValidMailDomains) == 0 {
		return true
	}

	for _, domain := range settings.Current.ValidMailDomains {
		if strings.HasSuffix(email, "@"+domain) {
			return true
		}
	}
	return false
}
