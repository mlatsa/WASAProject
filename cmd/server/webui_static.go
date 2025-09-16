package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func (rt *Router) registerAPI() {
	// Liveness
	rt.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Simple demo API
	rt.mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
	})
}

func (rt *Router) registerWebUI() {
	// Serve built frontend from webui/dist
	static := http.FileServer(http.Dir("webui/dist"))
	rt.mux.Handle("/assets/", static)
	rt.mux.Handle("/favicon.ico", static)

	// SPA fallback: serve index.html for routes that aren't real files
	rt.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := "webui/dist" + r.URL.Path
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			http.ServeFile(w, r, p)
			return
		}
		http.ServeFile(w, r, "webui/dist/index.html")
	})
}
