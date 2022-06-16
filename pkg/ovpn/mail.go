package ovpn

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-mail/mail"

	"fhyx.tech/platform/ovpntend/ui"
)

// vars
var (
	ErrInvalidName  = errors.New("invalid name")
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
	if !IsValidName(name) {
		return ErrInvalidName
	}
	body, err := GetClientConfig(ctx, name)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return ErrEmptyConfig
	}
	logger().Infow("read client body", "bytes", len(body))

	tpl, err := ui.Load("mail/" + oscat + ".htm")
	if err != nil {
		logger().Infow("mail data empty", "name", name, "osc", oscat, "err", err)
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
