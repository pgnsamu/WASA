package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/users/:id/conversations", rt.newConversation)
	rt.router.POST("/users/:id/photo", rt.setMyPhoto)
	rt.router.PUT("/users/:id/username", rt.setMyUserName)
	rt.router.GET("/users/:id", rt.getUserInfo)
	rt.router.GET("/users", rt.getUsers)
	rt.router.GET("/users/:id/conversations/:conversationId/users", rt.getParticipants)
	rt.router.DELETE("/users/:id/conversations/:conversationId/users/:toDelete", rt.leaveGroup)
	rt.router.POST("/users/:id/conversations/:conversationId/users", rt.addToGroup)
	rt.router.GET("/users/:id/conversations", rt.getMyConversations)
	rt.router.GET("/users/:id/conversations/:conversationId", rt.GetConversationInfoReq)
	rt.router.POST("/users/:id/conversations/:conversationId/group", rt.setGroupName) // TODO: rinominare endpoint?
	rt.router.POST("/users/:id/conversations/:conversationId/photo", rt.setGroupPhoto)
	rt.router.POST("/users/:id/conversations/:conversationId/messages", rt.sendMessageReq)
	rt.router.POST("/users/:id/conversations/:conversationId/messages/:messageId", rt.postForwardMessage)
	rt.router.DELETE("/users/:id/conversations/:conversationId/messages/:messageId", rt.deleteMessage)

	rt.router.POST("/session", rt.doLogin)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
