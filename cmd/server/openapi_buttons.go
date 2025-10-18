package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Used by PUT /api/groups/{id}/name
type groupNameBody struct {
	Name string `json:"name"`
}

// IMPORTANT: call r.registerButtonAPI() BEFORE the catch-all.
func (rt *Router) registerButtonAPI() {
	// /api/healthz -> same as /healthz, for frontend proxy convenience
	rt.mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// POST /api/session -> { identifier: "..." }
	rt.mux.HandleFunc("/api/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"identifier": "bearer-demo-token"})
	})

	// PUT /api/user/username -> { updated: true }
	rt.mux.HandleFunc("/api/user/username", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"updated": true})
	})

	// PUT /api/user/photo -> 204
	rt.mux.HandleFunc("/api/user/photo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// Group ops: /api/groups/{id}/members|leave|name|photo
	rt.mux.HandleFunc("/api/groups/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/groups/") {
			http.NotFound(w, r)
			return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/groups/")
		parts := strings.Split(rest, "/") // expect exactly 2: [id, action]
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		id, action := parts[0], parts[1]
		if id == "" {
			http.NotFound(w, r)
			return
		}

		switch action {
		case "members":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"added": true})

		case "leave":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case "name":
			if r.Method != http.MethodPut {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			defer r.Body.Close()
			var body groupNameBody
			_ = json.NewDecoder(r.Body).Decode(&body) // simple stub; ignore decode errors
			if strings.TrimSpace(body.Name) == "" {
				http.Error(w, "name required", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"renamed": true})

		case "photo":
			if r.Method != http.MethodPut {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.NotFound(w, r)
		}
	})
}
