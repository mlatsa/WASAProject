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

	// Messaging API routes
	r.router.GET("/", r.getHelloWorld)
	r.router.GET("/liveness", r.liveness)

	r.router.POST("/session", r.doLogin)

	r.router.GET("/conversations", r.getMyConversations)
	r.router.GET("/conversations/:id", r.getConversation)
	r.router.POST("/conversations/:id/messages", r.sendMessage)

	return r
}

func (rt *Router) Handler() http.Handler {
	return rt.router
}
