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
	h := withCORS(rt.Handler())
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}
