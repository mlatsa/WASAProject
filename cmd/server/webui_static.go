package main

import (
	"net/http"
	"os"
	"path/filepath"
)

// Serves files from webui/dist and falls back to index.html for SPA routes.
func (rt *Router) registerWebUI() {
	dist := "webui/dist"
	rt.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(dist, "index.html"))
			return
		}
		path := filepath.Join(dist, r.URL.Path)
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			http.ServeFile(w, r, path)
			return
		}
		http.ServeFile(w, r, filepath.Join(dist, "index.html"))
	})
}
