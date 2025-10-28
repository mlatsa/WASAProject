package main

import (
	"log"
	"net/http"
	"os"
)

type Router struct{ mux *http.ServeMux }

type Config struct{ Port string }

func loadConfiguration() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{Port: port}
}

// Base liveness probe at /healthz
func (rt *Router) registerAPI() {
	rt.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func main() {
	cfg := loadConfiguration()

	r := &Router{mux: http.NewServeMux()}

	// ---- ORDER MATTERS: specific -> catch-alls -> static UI ----
	r.registerAPI()             // /healthz
	r.registerButtonAPI()       // /api/session, /api/user/*, /api/groups/*
	r.registerChatAPI()         // /api/conversations/*, /api/messages/*
	r.registerOpenAPIImpl()     // /openapi.json, /openapi.yaml (if present)
	r.registerWebUI()           // static SPA (/) LAST

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r.mux))
}
