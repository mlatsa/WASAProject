package api

import (
	"net/http"
)

// Handler returns an instance of httprouter. Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/user", rt.setMyUserName)
	rt.router.GET("/users", rt.getAllUsers)
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.GET("/conversations/:to_user", rt.getConversation)
	rt.router.POST("/messages", rt.sendMessage)
	rt.router.DELETE("/messages/:msg_id", rt.deleteMessage)
	rt.router.POST("/messages/:msg_id", rt.forwardMessage)
	rt.router.PUT("/user/photo", rt.setMyPhoto)
	rt.router.POST("/messages/:msg_id/comments", rt.commentMessage)
	rt.router.DELETE("/messages/:msg_id/comments", rt.uncommentMessage)
	rt.router.GET("/groups", rt.getMyGroups)
	rt.router.PUT("/groups/:id", rt.addToGroup)
	rt.router.DELETE("/groups/:id", rt.leaveGroup)
	rt.router.POST("/groups/:id", rt.setGroupName)
	rt.router.GET("/groups/:id", rt.getMembers)
	rt.router.PUT("/group/photo", rt.setGroupPhoto)

	return rt.router
}
