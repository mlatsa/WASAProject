package main

import (
	"log"
	"net/http"

	"github.com/mlatsa/WASAProject/service/api"
)

func main() {
	rt := api.NewRouter()
	addr := ":3000"
	log.Printf("listening on %s\n", addr)
	if err := http.ListenAndServe(addr, rt.Handler()); err != nil {
		log.Fatal(err)
	}
}
