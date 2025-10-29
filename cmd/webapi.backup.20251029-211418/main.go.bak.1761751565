package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mlatsa/WASAProject/service/api"
	"github.com/mlatsa/WASAProject/service/extproxy"
)

func main() {
	rt := api.NewRouter()

	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// External proxy, only if UPSTREAM_API is set
	if os.Getenv("UPSTREAM_API") != "" {
		mux.Handle("/ext/", extproxy.MustProxy("/ext/"))
	}

	// Your API router
	mux.Handle("/", rt.Handler())

	addr := ":3000"
	log.Printf("listening on %s\n", addr)
	h := withCORS(mux)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}
