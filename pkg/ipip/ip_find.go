package ipip

import (
	"log"
	"os"

	"github.com/ipipdotnet/datx-go"
)

var cityP *datx.City

func init() {
	dir := os.Getenv("IPIP_DATX_PATH")
	if dir == "" {
		log.Print("IPIP_DATX_PATH not found")
		return
	}
	var err error
	cityP, err = datx.NewCity(dir)
	if err != nil {
		log.Printf("load ipip datx file failed, ERR %s", err)
		return
	}
}

// FindCity ...
func FindCity(ip string) (city, province, country string) {
	if cityP == nil {
		return
	}
	arr, err := cityP.Find(ip)
	if err != nil {
		log.Printf("find ip %s ERR %s", ip, err)
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
