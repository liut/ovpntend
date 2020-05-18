package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"fhyx.tech/platform/ovpntend/pkg/assets"
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

	if debug {
		s.ar.NotFound(http.FileServer(http.Dir("./ui")).ServeHTTP)
	} else {
		s.ar.NotFound(assets.ServeHTTP)
		// s.ar.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// 	if strings.HasPrefix(r.RequestURI, "/static") {
		// 		assets.ServeHTTP(w, r)
		// 		return
		// 	}
		// 	logger().Infow("not found", "uri", r.RequestURI)
		// 	http.NotFound(w, r)
		// })
	}

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
			logger().Infow("router walk fail", "err", err)
		}
	}
	return s
}

func (s *server) Serve() {
	logger().Infow("listen web server", "addr", s.Addr)
	err := s.hs.ListenAndServe()
	if err != nil {
		logger().Errorw("server start fail", "err", err)
	}
}

func (s *server) Stop() {
	if s.hs != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		s.hs.Shutdown(ctx)
	}
}
