package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/users/:id/conversations", rt.postConversation)
	rt.router.POST("/users/:id/photo", rt.postProfileImage)
	rt.router.PUT("/users/:id/username", rt.putUsername)
	rt.router.GET("/users/:id", rt.getUserInfo)
	rt.router.GET("/users", rt.getUsers)
	rt.router.GET("/users/:id/conversations/:conversationId/users", rt.getParticipants)
	rt.router.DELETE("/users/:id/conversations/:conversationId/users/:toDelete", rt.delParticipant)
	rt.router.POST("/session", rt.createUser)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
