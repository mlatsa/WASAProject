package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router *httprouter.Router
	store  *Store
}

func NewRouter() *Router {
	rt := &Router{
		router: httprouter.New(),
		store:  newStore(),
	}
	rt.registerRoutes()
	return rt
}

func (rt *Router) Handler() http.Handler { return rt.router }
