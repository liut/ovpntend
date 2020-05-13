package web

import (
	"net/http"

	"github.com/go-chi/chi"

	staffio "github.com/liut/staffio-client"
)

// User online user
type User = staffio.User

// vars from staffio
var (
	SetLoginPath    = staffio.SetLoginPath
	SetAdminPath    = staffio.SetAdminPath
	UserFromContext = staffio.UserFromContext
)

func (s *server) strapRouter() {

	s.ar.Get("/", handleNoContent)
	s.ar.Get("/ping", handlerPing)

	s.ar.Route("/auth", func(r chi.Router) {
		r.Get("/login", staffio.LoginHandler)
		r.Get("/logout", staffio.LogoutHandler)
		r.Method(http.MethodGet, "/callback", staffio.AuthCodeCallback())
	})

	s.ar.Route("/api/vpn", func(r chi.Router) {
		r.Use(staffio.Middleware(staffio.WithRefresh()))
		r.Get("/names", handlerNames)
		r.Get("/status/{idx}", handlerStatus)
		r.Post("/client/send", handlerSendClient)
	})

	// TODO: show html templates

}

func handleNoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong\n"))
}
