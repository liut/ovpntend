package main

import (
	"flag"
	"log"
	"net/http"
	// "os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	envcfg "github.com/wealthworks/envflagset"

	"fhyx.tech/platform/openvpn-monitor/ovpn/status"
)

var (
	version = "dev"
	dir     string
	listen  string
	ovpnmgr string
)

func init() {
	envcfg.New("ovpn", version)
	flag.StringVar(&listen, "listen", ":7605", "listen addr for http")
	flag.StringVar(&dir, "status-dir", "", "")
	flag.StringVar(&ovpnmgr, "man-addr", "127.0.0.1:7505", "management addr")
}

func main() {
	envcfg.Parse()
	// if dir == "" {
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	router := chi.NewRouter()
	router.Get("/ovpn/status", handlerStatus)

	log.Printf("listen http server on %s", listen)
	log.Print(http.ListenAndServe(listen, router))
}

func handlerStatus(w http.ResponseWriter, req *http.Request) {
	var result *status.Status
	var err error
	if dir != "" {
		filename := filepath.Join(dir, "status.log")
		log.Print(filename)
		result, err = status.ParseFile(filename)
	} else {
		log.Printf("read from %s", ovpnmgr)
		result, err = status.ParseAddr(ovpnmgr)
	}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	render.JSON(w, req, result)
}
