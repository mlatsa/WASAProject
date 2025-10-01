package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (rt *Router) registerOpenAPICatchAll() {
	registerCatchAllWithPrefixes(rt, []string{
		"/api/",
		"/v1/",
	})
}

func registerCatchAllWithPrefixes(rt *Router, prefixes []string) {
	for _, p := range prefixes {
		prefix := p // capture
		rt.mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, prefix) {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "created"})
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
			}
		})
	}
}
