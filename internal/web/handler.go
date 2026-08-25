package web

import (
	"embed"
	"io/fs"
	"net/http"
)

func SPA(dist embed.FS) http.Handler {
	sub, err := fs.Sub(dist, "frontend/dist")
	if err != nil {
		panic(err)
	}

	files := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Stat(sub, r.URL.Path[1:]); err != nil {
			http.ServeFileFS(w, r, sub, "index.html") // SPA fallback
			return
		}
		files.ServeHTTP(w, r)
	})
}
