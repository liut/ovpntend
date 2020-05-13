package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"fhyx.tech/gopak/binding"

	"fhyx.tech/platform/ovpntend/pkg/ovpn"
	"fhyx.tech/platform/ovpntend/pkg/settings"
	"fhyx.tech/platform/ovpntend/pkg/status"
)

func apiError(w http.ResponseWriter, r *http.Request, status int, err interface{}) {
	res := render.M{
		"status": status,
		"error":  err,
	}
	switch ret := err.(type) {
	case error:
		res["message"] = ret.Error()
	case fmt.Stringer:
		res["message"] = ret.String()
	case string, *string, []byte:
		res["message"] = ret
	}
	render.JSON(w, r, res)
}

func apiOk(w http.ResponseWriter, r *http.Request, data interface{}, count int) {
	res := render.M{"status": 0}
	if data != nil {
		res["data"] = data
	}
	if count > 0 {
		res["count"] = count
	}
	render.JSON(w, r, res)
}

func handlerNames(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{"status": 0, "names": settings.Current.ManageNames})
}

func handlerStatus(w http.ResponseWriter, r *http.Request) {
	count := len(settings.Current.ManageAddrs)
	if count == 0 {
		w.WriteHeader(204)
		return
	}
	var idx int
	if s := chi.URLParam(r, "idx"); s != "" {
		idx, _ = strconv.Atoi(s)
	}
	if idx >= count {
		w.WriteHeader(400)
		return
	}
	var result *status.Status
	var err error

	ovpnmgr := settings.Current.ManageAddrs[idx]
	logger().Infow("read ovpn status", "addr", ovpnmgr)
	result, err = status.ParseAddr(ovpnmgr)

	if err != nil {
		logger().Infow("read fail", "err", err)
		http.Error(w, err.Error(), 400)
		return
	}
	if len(settings.Current.ManageNames) > idx {
		result.Label = settings.Current.ManageNames[idx]
	}
	render.JSON(w, r, render.M{"status": 0, "clients": result.ClientList, "name": result.Label})
}

type getClientParam struct {
	Name  string `json:"name"`
	OSCat string `json:"oscat"`
}

func handlerSendClient(w http.ResponseWriter, r *http.Request) {
	param := new(getClientParam)
	if err := binding.Bind(r, param); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := ovpn.SendConfig(r.Context(), param.Name, param.OSCat); err != nil {
		logger().Infow("send fail", "err", err)
		http.Error(w, err.Error(), 503)
		return
	}

	apiOk(w, r, nil, 0)
}
