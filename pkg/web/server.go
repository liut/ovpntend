package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/liut/ovpntend/ui"
)

type server struct {
	Addr string

	ar *chi.Mux     // app router
	hs *http.Server // http server
}

// New return new web server
func New(debug bool, addr string) interface {
	Serve()
	Stop()
} {
	inDev = debug
	ar := chi.NewMux()
	if debug {
		ar.Use(middleware.Logger)
	}
	ar.Use(middleware.Recoverer)

	s := &server{Addr: addr, ar: ar}
	s.strapRouter()

	s.ar.NotFound(ui.Handler().ServeHTTP)

	s.hs = &http.Server{
		Addr:    s.Addr,
		Handler: s.ar,
	}

	if debug {
		walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			route = strings.Replace(route, "/*/", "/", -1)
			fmt.Printf("DEBUG: %-6s %-24s --> %s (%d mw)\n", method, route, nameOfFunction(handler), len(middlewares))
			return nil
		}

		if err := chi.Walk(ar, walkFunc); err != nil {
			slog.Info("router walk fail", "err", err)
		}
	}
	return s
}

func (s *server) Serve() {
	slog.Info("listen web server", "addr", s.Addr)
	err := s.hs.ListenAndServe()
	if err != nil {
		slog.Error("server start fail", "err", err)
	}
}

func (s *server) Stop() {
	if s.hs != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		s.hs.Shutdown(ctx)
	}
}
