package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/mlatsa/WASAProject/service/api"
)

func main() {
	rt := api.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:5173", // Vite dev
			"http://localhost:8080", // Nginx default
			"http://localhost:8081", // Alt Nginx port you used
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	addr := ":3000"
	log.Printf("listening on %s\n", addr)
	if err := http.ListenAndServe(addr, cors(rt.Handler())); err != nil {
		log.Fatal(err)
	}
}
