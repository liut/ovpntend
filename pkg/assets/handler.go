package assets

import (
	"net/http"
)

//go:generate staticfiles --package kcstatic -o files.go ../../ui/static

// Handler all static/ files embedded as a Go library
func Handler(inDev bool) http.Handler {
	if inDev {
		return http.FileServer(http.Dir("ui/static/"))
	}

	return Server
}
