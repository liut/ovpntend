package ovpn

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	expect "github.com/ThomasRooney/gexpect"

	"fhyx.tech/platform/ovpntend/pkg/settings"
)

var (
	caPassword string
	timeout    = 15

	cmdEnv []string
)

func init() {
	caPassword = settings.Current.EasyRSACaPassword
	cmdEnv = append(os.Environ(),
		"EASYRSA="+settings.Current.EasyRSABIN,
		"EASYRSA_PKI="+settings.Current.EasyRSAPKI,
		"OVPN_CN="+settings.Current.OpenVPNHost,
		"OVPN_DEFROUTE=0",
		fmt.Sprintf("OVPN_PORT=%d", settings.Current.OpenVPNPort),
		"OVPN_PROTO="+settings.Current.OpenVPNProto,
	)
}

// CreateCertificate ...
func CreateCertificate(name string) {
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
		"host":  settings.Current.OpenVPNHost,
		"port":  settings.Current.OpenVPNPort,
		"proto": settings.Current.OpenVPNProto,
		"dev":   "tun",
	}
	names := map[string]string{
		"key":  "private/" + name + ".key",
		"cert": "issued/" + name + ".crt",
		"ca":   "ca.crt",
		"dh":   "dh.pem",
		"ta":   "ta.key",
	}

	dir := settings.Current.EasyRSAPKI

	for k, file := range names {
		var b []byte
		b, err = ioutil.ReadFile(path.Join(dir, file))
		if err != nil {
			logger().Infow("read fail", "file", file, "err", err)
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
