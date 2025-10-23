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

func (rt *Router) registerAPI() {
	// base liveness on /healthz (mirrored by /api/healthz in button API)
	rt.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func main() {
	cfg := loadConfiguration()

	r := &Router{mux: http.NewServeMux()}
	r.registerAPI()
	r.registerButtonAPI() // /api/session, /api/user/*, /api/groups/*
	r.registerChatAPI()   // /api/conversations/*, /api/messages/*
	r.registerWebUI()     // serve built SPA from webui/dist

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r.mux))
}
