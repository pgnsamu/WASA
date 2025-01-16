package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.GET("/users", rt.getUsers)                                                               // non sta nel doc
	rt.router.GET("/users/:id/conversations/:conversationId/users", rt.getParticipants)                // non sta nel doc
	rt.router.POST("/info", rt.getLoggedUserInfo)                                                      // non sta nel doc e non viene usata (per ora)
	rt.router.GET("/users/:id/conversations/:conversationId/messages", rt.getMessagesFromConversation) // non sta nel doc

	rt.router.POST("/session", rt.doLogin)                                                                                    // sta nel doc // TODO: da cambiare nel doc
	rt.router.GET("/users/:id", rt.getUserInfo)                                                                               // sta nel doc
	rt.router.PUT("/users/:id/username", rt.setMyUserName)                                                                    // TODO: togliere put e mettere post sta nel doc
	rt.router.POST("/users/:id/photo", rt.setMyPhoto)                                                                         // sta nel doc
	rt.router.POST("/users/:id/conversations/:conversationId/users", rt.addToGroup)                                           // sta nel doc
	rt.router.DELETE("/users/:id/conversations/:conversationId/users", rt.leaveGroup)                                         // TODO: nel doc è scritto senza toDelete perché era pensato in modo che un utente non possa espellere altri
	rt.router.PUT("/users/:id/conversations/:conversationId/group", rt.setGroupName)                                          // STA NEL DOC
	rt.router.PUT("/users/:id/conversations/:conversationId/photo", rt.setGroupPhoto)                                         // sta nel doc
	rt.router.GET("/users/:id/conversations", rt.getMyConversations)                                                          // sta nel doc
	rt.router.POST("/users/:id/conversations", rt.newConversation)                                                            // sta nel doc
	rt.router.GET("/users/:id/conversations/:conversationId", rt.getConversation)                                             // sta nel doc
	rt.router.POST("/users/:id/conversations/:conversationId/messages", rt.sendMessage)                                       // sta nel doc
	rt.router.POST("/users/:id/conversations/:conversationId/messages/:messageId", rt.forwardMessage)                         // sta nel doc
	rt.router.DELETE("/users/:id/conversations/:conversationId/messages/:messageId", rt.deleteMessage)                        // sta nel doc
	rt.router.POST("/users/:id/conversations/:conversationId/messages/:messageId/comments", rt.commentMessage)                // sta nel doc
	rt.router.DELETE("/users/:id/conversations/:conversationId/messages/:messageId/comments/:commentId", rt.uncommentMessage) // sta nel doc

	// Special routes
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	return rt.router
}
