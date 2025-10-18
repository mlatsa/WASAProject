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

func main() {
	cfg := loadConfiguration()

	r := &Router{mux: http.NewServeMux()}

	// Register concrete APIs first
	r.registerAPI()       // /healthz, /api/healthz, /api/version
	r.registerButtonAPI() // /api/session, /api/user/*, /api/groups/*
	r.registerChatAPI()   // /api/conversations/*, /api/messages/*

	// Keep generic/catch-all AFTER specific handlers if you have them:
	// r.registerOpenAPIImpl()
	// r.registerOpenAPICatchAll()

	// Static Web UI (serves built SPA from webui/dist)
	r.registerWebUI()

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r.mux))
}
