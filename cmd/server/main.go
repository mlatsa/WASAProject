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
	r.registerAPI()
	r.registerOpenAPIImpl()
        r.registerButtonAPI()
        r.registerChatAPI()
	r.registerOpenAPICatchAll()
	r.registerWebUI()

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r.mux))
}
