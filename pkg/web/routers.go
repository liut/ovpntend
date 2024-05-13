package web

import (
	"net/http"

	"github.com/go-chi/chi"

	staffio "github.com/liut/staffio-client"

	"github.com/liut/ovpntend/pkg/settings"
)

// User online user
type User = staffio.User

// vars from staffio
var (
	SetLoginPath    = staffio.SetLoginPath
	SetAdminPath    = staffio.SetAdminPath
	UserFromContext = staffio.UserFromContext

	authzr staffio.Authorizer
)

func init() {
	authzr = staffio.NewAuth(staffio.WithRefresh(), staffio.WithURI(staffio.LoginPath), staffio.WithCookie(
		settings.Current.CookieName,
		settings.Current.CookiePath,
		settings.Current.CookieDomain,
	))
}

func (s *server) strapRouter() {

	var suffix string
	if settings.Current.ServerPlace != "" {
		suffix = "/" + settings.Current.ServerPlace
	}

	s.ar.Get("/ping", handlerPing)

	s.ar.Route("/auth", func(r chi.Router) {
		r.Get("/login", staffio.LoginHandler)
		r.Get("/logout", staffio.LogoutHandler)
		r.Method(http.MethodGet, "/callback", staffio.AuthCodeCallback())
	})

	s.ar.Route("/api/vpn"+suffix, func(r chi.Router) {
		r.Use(authzr.Middleware())
		r.Get("/names", handlerNames)
		r.Get("/status/{idx}", handlerStatus)
		r.Post("/client/send", handlerSendClient)
	})

	// s.ar.Get("/", handleNoContent)
	SetAdminPath("/")
	s.ar.Group(func(r chi.Router) {
		r.Use(authzr.MiddlewareWordy(true))
		r.Get("/", handlerHome)
		r.Get("/status{idx}", handlerStatus)
		r.Get("/status", handlerStatus)
		r.Get("/route-customize", handleRoutes)
	})

}

// nolint
func handleNoContent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(204)
}

// nolint
func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong\n"))
}
