package webui

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler can be used if code expects a webui.Handler http.Handler.
var Handler http.Handler = http.NotFoundHandler()

// AddRoutes is a common pattern: if your server calls webui.AddRoutes(router), this no-ops.
func AddRoutes(r *httprouter.Router) {}

// Mount allows code that expects webui.Mount(mux, prefix) to compile; it serves 404s by default.
func MountMux(mux *http.ServeMux, prefix string) {}

// MountRouter allows code that expects webui.Mount(router) to compile; it no-ops.
func MountRouter(r *httprouter.Router) {}
