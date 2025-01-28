package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Session and user management
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/user/username", rt.authenticate(rt.changeUserName))
	rt.router.PUT("/user/photo", rt.authenticate(rt.changeUserPhoto))

	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// user getting and existence verification.
	rt.router.GET("/users", rt.authenticate(rt.getAllUsers))
	rt.router.GET("/users/photo/:username", rt.authenticate(rt.getUserPhoto))

	rt.router.GET("/conversations", rt.authenticate(rt.getUserConversations))
	// route protette dal membership check per assicurarsi che l'utente faccia parte della conversazione.
	rt.router.GET("/conversations/info/:conversationId", rt.authenticate(rt.checkUserInConversation(rt.GetConversationInfo)))
	rt.router.POST("/conversations/create", rt.authenticate(rt.CreateConversation))
	rt.router.PUT("/conversations/adduser", rt.authenticate(rt.checkUserInConversation(rt.AddParticipant)))
	rt.router.PUT("/conversations/kickuser", rt.authenticate(rt.checkUserInConversation(rt.RemoveParticipant)))

	rt.router.POST("/conversations/message", rt.authenticate(rt.checkUserInConversation(rt.SendMessage)))
	rt.router.GET("/conversations/message/info/:messageId", rt.authenticate(rt.checkUserInConversation(rt.GetMessageInfo)))
	rt.router.POST("/conversations/message/forward", rt.authenticate(rt.checkUserInConversation((rt.ForwardMessage))))
	rt.router.DELETE("/conversations/message/delete/:messageId", rt.authenticate(rt.checkUserInConversation(rt.DeleteMessage)))
	rt.router.PUT("/conversations/photo", rt.authenticate(rt.checkUserInConversation(rt.UpdateConversationPhoto)))
	rt.router.PUT("/conversations/name", rt.authenticate(rt.checkUserInConversation(rt.UpdateConversationName)))

	// comments
	// They all include the conversation membership check, since comments are tied to messages, which are tied to conversations.
	rt.router.POST("/conversations/message/comment/:messageId", rt.authenticate(rt.checkUserInConversation(rt.AddComment)))
	rt.router.DELETE("/conversations/message/comment/:commentId", rt.authenticate(rt.checkUserInConversation(rt.DeleteComment)))
	rt.router.GET("/conversations/message/comment/:commentId", rt.authenticate(rt.checkUserInConversation(rt.GetCommentInfo)))
	rt.router.GET("/conversations/message/getcomments/:messageId", rt.authenticate(rt.checkUserInConversation(rt.GetComments)))

	return rt.router
}
