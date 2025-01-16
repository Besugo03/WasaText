package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	Name string `json:"name"`
}

type loginResponse struct {
	Identifier string `json:"identifier"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the user credentials (this is just a placeholder, replace with actual validation logic)
	if req.Name != "testuser" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a unique identifier for the session (this is just a placeholder, replace with actual logic)
	// TODO implement a secure way to generate the identifier. I'll probably use a UUID or just append something to the username.
	identifier := "abcdef012345"

	resp := loginResponse{
		Identifier: identifier,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extract the token from the Authorization header
		token := r.Header.Get("Authorization")
		rt.baseLogger.Infof("Received token:", token)

		// Check if the token is present
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// print the received token
		// Validate the token (this is just a placeholder, replace with actual validation logic)
		// TODO implement a secure way to validate the token. I'll probably stick them in a database.
		if token != "abcdef012345" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r, ps)
	}
}
