package ovpn

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	expect "github.com/ThomasRooney/gexpect"
)

// CreateCertificate ...
func CreateCertificate(name string) {
	// TODO: set env for cmd
	if len(name) == 0 {
		logger().Fatalw("empty name")
		return
	}
	e, err := expect.Spawn("easyrsa build-client-full " + name + " nopass")
	if err != nil {
		logger().Fatalw("call easyrsa fail", "name", name, "err", err)
	}
	defer e.Close()

	e.Expect("/private/ca.key:")
	e.Send(caPassword + "\n")
}

// GenerateCRL ...
func GenerateCRL() {
	// TODO: set env for cmd
	e, err := expect.Spawn("easyrsa gen-crl")
	if err != nil {
		logger().Fatalw("call easyrsa fail", "err", err)
	}
	e.Expect("/private/ca.key:")
	e.Send(caPassword + "\n")
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
		b, err = ioutil.ReadFile(path.Join(easyrasPKI, file))
		if err != nil {
			logger().Infow("read fail", "file", file, "err", err)
			err = fmt.Errorf("file %q not found", file)
			return
		}
		logger().Debugw("read ok", "file", file)
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
