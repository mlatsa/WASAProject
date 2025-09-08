package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router *httprouter.Router
}

func NewRouter() *Router {
	r := &Router{router: httprouter.New()}
	r.registerRoutes()
	return r
}

func (rt *Router) Handler() http.Handler {
	return rt.router
}
