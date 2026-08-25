package main

import (
	"embed"
	"net/http"

	"onlystudents/internal/web"

	"github.com/go-chi/chi/v5"
)

//go:embed all:frontend/dist
var dist embed.FS

func main() {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong!"))
	})
	r.Get("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"hello world"}`))
	})
	r.Handle("/*", web.SPA(dist))
	http.ListenAndServe(":8080", r)
}
