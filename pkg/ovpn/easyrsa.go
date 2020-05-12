package ovpn

import (
	"context"
	"fmt"
	"os"
	"os/exec"

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
	cmd := exec.CommandContext(ctx, "sh", "script/ovpn_getclient", name)
	cmd.Env = cmdEnv
	out, err = cmd.Output()
	if err != nil {
		logger().Warnw("call ovpn_getclient fail", "name", name, "err", err)
	}
	return
}
