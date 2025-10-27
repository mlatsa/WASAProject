package api

func (rt *Router) registerRoutes() {
	// health
	rt.router.GET("/health", rt.health)
	rt.router.GET("/liveness", rt.health) // legacy OK

	// auth
	rt.router.POST("/session", rt.doLogin)

	// user profile
	rt.router.PUT("/user/username", rt.putUserUsername)
	rt.router.PUT("/user/photo", rt.putUserPhoto)

	// conversations
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.GET("/conversations/:conversationId", rt.getConversation)
	rt.router.POST("/conversations/:conversationId/messages", rt.sendMessage)

	// messages – forward & reactions & delete
	rt.router.POST("/messages/:messageId/forward", rt.postMessageForward)
	rt.router.POST("/messages/:messageId/reactions", rt.postMessageReaction)
	rt.router.DELETE("/messages/:messageId/reactions/:reactionId", rt.deleteMessageReaction)
	rt.router.DELETE("/messages/:messageId", rt.deleteMessage)

	// groups (conversation-scoped)
	rt.router.POST("/groups/:conversationId/members", rt.postGroupMember)
	rt.router.POST("/groups/:conversationId/leave", rt.postGroupLeave)
	rt.router.PUT("/groups/:conversationId/name", rt.putGroupName)
	rt.router.PUT("/groups/:conversationId/photo", rt.putGroupPhoto)
}
