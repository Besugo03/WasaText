package api

// Contains all the needed functions to correctly log in and authenticate the user with each
// protected request.

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	Username string `json:"username"`
}

type loginResponse struct {
	Identifier string `json:"identifier"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user_created = false
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// query the database to see if the identifier exists
	identifier, err := rt.db.GetUser(req.Username)
	if err != nil {
		// if it's a no rows error, it means the user doesn't exist, and needs to be created.
		if err.Error() == "sql: no rows in result set" {
			rt.baseLogger.Infof("no rows for %v exist in database. Creating an entry...", req.Username)
			// Generate a unique identifier for the session (this is just a placeholder, replace with actual logic)
			// TODO implement a way to generate the identifier based on the name.
			// I'll probably use a UUID or just append something to the username.
			identifier = req.Username + "123"
			err := rt.db.SetUser(req.Username, identifier)
			if err != nil {
				rt.baseLogger.WithError(err).Errorf("error creating user %v", req.Username)
				errorstring := err.Error()
				http.Error(w, "Failed to create user : "+errorstring, http.StatusInternalServerError)
				return
			}
			rt.baseLogger.Infof("User %v created successfully! (id : %v)", req.Username, identifier)
			user_created = true

		} else {
			// if it's not a no rows error, it's a database error.
			rt.baseLogger.Infof("Error querying the database for user %v (%v): %v", identifier, req.Username, err)
			http.Error(w, "Failed to query the database", http.StatusInternalServerError)
			return
		}
	}

	resp := loginResponse{
		Identifier: identifier,
	}

	w.Header().Set("Content-Type", "application/json")
	if user_created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// authenticate is a middleware that checks if the request is authenticated.
// It checks if the request contains a valid token in the Authorization header.
// If the token is valid, it calls the next handler, otherwise it returns an Unauthorized error.
func (rt *_router) authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extract the token from the Authorization header
		token := r.Header.Get("Authorization")
		rt.baseLogger.Infof("Received token:", token)
		// remove the "Bearer "
		token = strings.TrimPrefix(token, "Bearer ")

		// Check if the token is present
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the token matches any user in the database
		_, err := rt.db.GetUserName(token)
		if err != nil {
			rt.baseLogger.WithError(err).Errorf("Error querying the database for token %v", token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		rt.baseLogger.Infof("Token %v is valid", token)

		// Call the next handler if the token is valid.
		next(w, r, ps)
	}
}
