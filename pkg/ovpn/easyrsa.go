package ovpn

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"text/template"

	expect "github.com/agorer/gexpect"
)

// CreateCertificate ...
func CreateCertificate(name string) {
	// TODO: set env for cmd
	if len(name) == 0 {
		slog.Error("empty name")
		return
	}
	e, err := expect.Spawn("easyrsa build-client-full " + name + " nopass")
	if err != nil {
		slog.Error("call easyrsa fail", "name", name, "err", err)
		return
	}
	defer e.Close()

	if err := e.Expect("/private/ca.key:"); err != nil {
		slog.Info("expect fail", "err", err)
	}
	_ = e.Send(caPassword + "\n")
}

// GenerateCRL ...
func GenerateCRL() {
	// TODO: set env for cmd
	e, err := expect.Spawn("easyrsa gen-crl")
	if err != nil {
		slog.Error("call easyrsa fail", "err", err)
		return
	}
	if err := e.Expect("/private/ca.key:"); err != nil {
		slog.Info("expect fail", "err", err)
	}
	_ = e.Send(caPassword + "\n")
}

// GetClientConfig ...
func GetClientConfig(ctx context.Context, name string) (out []byte, err error) {
	if len(name) == 0 {
		err = ErrEmptyConfig
		return
	}
	cc := map[string]interface{}{
		"host":  serverHost,
		"port":  serverPort,
		"proto": serverProto,
		"dev":   "tun",
	}
	names := map[string]string{
		"cert": "issued/" + name + ".crt",
		"key":  "private/" + name + ".key",
		"ca":   "ca.crt",
		"dh":   "dh.pem",
		"ta":   "ta.key",
	}

	for k, file := range names {
		var b []byte
		b, err = os.ReadFile(path.Join(easyrasPKI, file))
		if err != nil {
			slog.Info("read fail", "file", file, "err", err)
			err = fmt.Errorf("file %q not found", file)
			return
		}
		slog.Debug("read ok", "file", file)
		cc[k] = string(bytes.TrimRight(b, "\n"))
	}

	var buf = new(bytes.Buffer)
	t := template.Must(template.New("cc").Parse(ccTpl))
	err = t.Execute(buf, cc)
	out = buf.Bytes()

	return
}

const (
	ccTpl = `
client
nobind
dev {{ .dev }}
remote-cert-tls server

remote {{ .host }} {{ .port }} {{ .proto }}

<key>
{{ .key }}
</key>
<cert>
{{ .cert }}
</cert>
<ca>
{{ .ca }}
</ca>
<dh>
{{ .dh }}
</dh>
<tls-auth>
{{ .ta }}
</tls-auth>
key-direction 1
`
)
