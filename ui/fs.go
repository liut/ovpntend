package ui

import (
	"embed"
	"net/http"
)

//go:embed mail/*.htm
//go:embed static/*.css
//go:embed static/*.js
//go:embed templates/*
var uifs embed.FS

func Load(name string) (string, error) {
	if data, err := uifs.ReadFile(name); err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func FS() http.FileSystem {
	return http.FS(uifs)
}
