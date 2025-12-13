package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// paths
	rt.router.POST("/session", rt.wrap(rt.postSession))
	rt.router.GET("/users", rt.wrap(rt.authHandler(rt.searchUser)))
	rt.router.PUT("/users/:userId/info/username", rt.wrap(rt.authHandler(rt.setMyUserName)))
	rt.router.PUT("/users/:userId/info/photo", rt.wrap(rt.authHandler(rt.setMyPhoto)))
	rt.router.GET("/users/:userId/conversations", rt.wrap(rt.authHandler(rt.getConversations)))
	rt.router.POST("/users/:userId/conversations", rt.wrap(rt.authHandler(rt.createConversation)))
	rt.router.GET("/users/:userId/conversations/:conversationId", rt.wrap(rt.authHandler(rt.getConversation)))
	rt.router.GET("/users/:userId/conversations/:conversationId/messages", rt.wrap(rt.authHandler(rt.getMessages)))
	rt.router.POST("/users/:userId/conversations/:conversationId/messages", rt.wrap(rt.authHandler(rt.postMessage)))
	rt.router.DELETE("/users/:userId/conversations/:conversationId/messages/:messageId", rt.wrap(rt.authHandler(rt.deleteMessage)))
	rt.router.POST("/users/{userId}/conversations/{conversationId}/messages/{MessageId}/forward_message", rt.wrap(rt.authHandler(rt.forwardMessage)))
	// FINIRE IL PROBLEMA DELLE CONVERSAZIONI BIDIREZIONALI
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
