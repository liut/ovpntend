package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	envcfg "github.com/wealthworks/envflagset"

	"fhyx/platform/openvpn-monitor/ovpn/status"
)

var (
	dir     string
	addr    string
	version = "dev"
)

func init() {
	envcfg.New("ovpn", version)
	flag.StringVar(&dir, "status-dir", "", "")
	flag.StringVar(&addr, "addr", ":8088", "listen addr for http")
}

func main() {
	envcfg.Parse()
	if dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Get("/status", handlerStatus)

	log.Printf("listen http server on %s", addr)
	http.ListenAndServe(addr, router)
}

func handlerStatus(w http.ResponseWriter, req *http.Request) {
	filename := filepath.Join(dir, "status.log")
	log.Print(filename)
	result, err := status.ParseFile(filename)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	render.JSON(w, req, result)
}
