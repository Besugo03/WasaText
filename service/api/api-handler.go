package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Session Route
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	// ? what the fuck does "wrap" even do?

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// test for protected route
	// quindi, per accedere a un endpoint protetto, é necessario wrapparlo nell'authenticate
	// (oltre che nel .wrap). Esempi di query presenti nel .py che uso per i test.
	rt.router.GET("/users/:user_id/conversations", rt.authenticate(rt.wrap(rt.getUserConversations)))

	return rt.router
}
