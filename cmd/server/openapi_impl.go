package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Add your real handlers here, one block per path in openapi.yaml.
func (rt *Router) registerOpenAPIImpl() {

	// ===== EXAMPLE: /api/v1/items (GET, POST) =====
	rt.mux.HandleFunc("/api/v1/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": "1", "name": "demo"},
			})
			return
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": "2"})
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ===== EXAMPLE: /api/v1/items/{id} (GET, PUT, DELETE) =====
	rt.mux.HandleFunc("/api/v1/items/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/items/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "name": "demo"})
			return
		case http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "updated": true})
			return
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// TODO: add blocks for each path in your openapi.yaml
}
