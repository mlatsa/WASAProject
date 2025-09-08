package webapi

import (
    "log"
    "net/http"

    "github.com/mlatsa/WASAProject/service/api"
)

// Main is the entrypoint used by the root main.go (for CI)
func Main() {
    rt := api.NewRouter()
    addr := ":3000"
    log.Printf("listening on %s\n", addr)
    if err := http.ListenAndServe(addr, rt.Handler()); err != nil {
        log.Fatal(err)
    }
}
