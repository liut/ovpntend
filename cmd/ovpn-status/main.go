package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	// "os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"fhyx.tech/platform/ovpntend/pkg/settings"
	"fhyx.tech/platform/ovpntend/pkg/status"
)

var (
	version   = "dev"
	dir       string
	listen    string
	showUsage bool
)

func init() {
	flag.StringVar(&listen, "listen", settings.Current.HTTPListen, "listen addr for http")
	flag.StringVar(&dir, "status-dir", settings.Current.StatusDir, "directory which has status.log of openvpn")
	flag.BoolVar(&showUsage, "usage", false, "show usage")
}

func main() {
	flag.Parse()
	if showUsage {
		settings.Usage()
		return
	}

	router := chi.NewRouter()
	router.Get("/ovpn/status/{idx}", handlerStatus)

	log.Printf("listen http server on %s", listen)
	log.Print(http.ListenAndServe(listen, router))
}

func handlerStatus(w http.ResponseWriter, req *http.Request) {
	count := len(settings.Current.ManageAddrs)
	if count == 0 {
		w.WriteHeader(204)
		return
	}
	var idx int
	if s := chi.URLParam(req, "idx"); s != "" {
		idx, _ = strconv.Atoi(s)
	}
	if idx >= count {
		w.WriteHeader(400)
		return
	}
	var result *status.Status
	var err error
	if dir != "" {
		filename := filepath.Join(dir, "status.log")
		log.Print(filename)
		result, err = status.ParseFile(filename)
	} else {
		ovpnmgr := settings.Current.ManageAddrs[idx]
		log.Printf("read from %s", ovpnmgr)
		result, err = status.ParseAddr(ovpnmgr)
	}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(settings.Current.ManageNames) >= idx+1 {
		result.Label = settings.Current.ManageNames[idx]
	}
	render.JSON(w, req, result)
}
