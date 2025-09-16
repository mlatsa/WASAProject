package main

import (
	"encoding/json"
	"net/http"
)

func (rt *Router) registerOpenAPICatchAll() {
	// Fallback for any /api/* that isn't explicitly handled elsewhere.
	rt.mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			// Many POST endpoints return 201
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "created",
			})
		case http.MethodDelete:
			// Many DELETE endpoints return 204 No Content
			w.WriteHeader(http.StatusNoContent)
		default:
			// GET, PUT, PATCH → return 200 OK with a tiny JSON
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
			})
		}
	})
}
