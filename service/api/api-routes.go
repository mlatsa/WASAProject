package api


func (rt *Router) registerRoutes() {
	r := rt.router

	r.GET("/health", rt.health)
	r.POST("/session", rt.doLogin)

	r.PUT("/user/username", rt.putUserUsername)
	r.PUT("/user/photo", rt.putUserPhoto)

	r.GET("/conversations", rt.getMyConversations)
	r.GET("/conversations/:conversationId", rt.getConversation)
	r.POST("/conversations/:conversationId/messages", rt.sendMessage)

	r.POST("/messages/:messageId/forward", rt.postMessageForward)
	r.POST("/messages/:messageId/reactions", rt.postMessageReaction)
	r.DELETE("/messages/:messageId/reactions/:reactionId", rt.deleteMessageReaction)
	r.DELETE("/messages/:messageId", rt.deleteMessage)

	// group stubs
	r.POST("/groups/:conversationId/members", rt.postGroupMember)
	r.POST("/groups/:conversationId/leave", rt.postGroupLeave)
	r.PUT("/groups/:conversationId/name", rt.putGroupName)
	r.PUT("/groups/:conversationId/photo", rt.putGroupPhoto)
}
