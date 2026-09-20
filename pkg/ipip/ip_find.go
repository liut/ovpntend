package ipip

import (
	"log/slog"
	"os"

	"github.com/ipipdotnet/datx-go"
)

var cityP *datx.City

func init() {
	dir := os.Getenv("IPIP_DATX_PATH")
	if dir == "" {
		slog.Warn("IPIP_DATX_PATH not found")
		return
	}
	var err error
	cityP, err = datx.NewCity(dir)
	if err != nil {
		slog.Warn("load ipip datx file failed", "err", err)
		return
	}
}

// FindCity ...
func FindCity(ip string) (city, province, country string) {
	if IsPrivate(ip) {
		city = "[内网]"
		return
	}
	if cityP == nil {
		return
	}
	arr, err := cityP.Find(ip)
	if err != nil {
		slog.Warn("find ip fail", "ip", ip, "err", err)
		return
	}

	country = arr[0]
	if len(arr) > 1 {
		province = arr[1]
		if len(arr) > 2 {
			city = arr[2]
		}
	}
	return
}
