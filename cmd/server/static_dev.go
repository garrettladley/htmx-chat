//go:build dev

package main

import (
	"net/http"
	"os"
)

func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("./cmd/server/public")))
}

func deps() http.Handler {
	return http.StripPrefix("/deps/", http.FileServer(http.Dir("./cmd/server/deps")))
}
