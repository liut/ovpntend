package ovpn

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-mail/mail"

	"fhyx.tech/platform/ovpntend/pkg/assets"
)

// vars
var (
	ErrEmptyConfig  = errors.New("empty config")
	ErrInvalidOS    = errors.New("invalid os")
	ErrMailTemplate = errors.New("mail template read fail")
)

// ParseOSCat ...
func ParseOSCat(osc string) bool {
	switch osc {
	case "linux", "mac", "windows":
		return true
	}
	return false
}

// SendConfig ...
func SendConfig(ctx context.Context, name, oscat string) error {
	body, err := GetClientConfig(ctx, name)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return ErrEmptyConfig
	}
	logger().Infow("read client body", "bytes", len(body))

	tpl := assets.Data("mail/" + oscat + ".htm")
	if 0 == len(tpl) {
		logger().Infow("mail data empty", "osc", oscat)
		return ErrMailTemplate
	}

	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", senderName, smtpUser))
	m.SetHeader("To", name)
	m.SetHeader("Subject", "Your OVPN Config!")
	m.SetBody("text/html", tpl)
	m.Attach(name+ssuffix+".ovpn", mail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(body)
		return err
	}))

	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		logger().Infow("send mail fail", "name", name, "err", err)
		return err
	}
	return nil
}
