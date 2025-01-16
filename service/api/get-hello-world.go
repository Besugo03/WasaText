package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain") // Set the content type to plain text
	_, _ = w.Write([]byte("Hello World!"))       // underscore is used to ignore the return value of Write.
	ping := rt.db.Ping()                         // Ping the database to check if it is still alive
	if ping != nil {                             // If the ping fails, log the error
		rt.baseLogger.WithError(ping).Error("error pinging the database")
	}
}
