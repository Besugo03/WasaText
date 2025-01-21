package api

// Contains all the needed functions to correctly log in and authenticate the user with each
// protected request.

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	identifier, err := rt.db.GetUserID(req.Username)
	if err != nil {
		// if it's a no rows error, it means the user doesn't exist, and needs to be created.
		if err.Error() == "sql: no rows in result set" {
			rt.baseLogger.Infof("no rows for %v exist in database. Creating an entry...", req.Username)
			// Generate a unique identifier for the session (this is just a placeholder, replace with actual logic)
			// TODO implement a way to generate the identifier based on the name.
			// I'll probably use a UUID or just append something to the username.
			identifier = req.Username + "123"
			err := rt.db.RegisterUser(req.Username, identifier, "")
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
		// rt.baseLogger.Infof("Received token:", token)
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
		// rt.baseLogger.Infof("Token %v is valid", token)

		// Call the next handler if the token is valid.
		next(w, r, ps)
	}
}

// checkUserInConversation is a middleware that checks if the user is a participant in the conversation, and if the conversation exists, before calling the next handler.
func (rt *_router) checkUserInConversation(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extract the token from the Authorization header
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		// rt.baseLogger.Infof("Received token : %v", token)

		var SpecifiedConversationID string

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

		// rt.baseLogger.Infof("Token %v exists", token)

		// Read the request body
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Reset the request body so it can be read again
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Now, we can proceed.

		// the request can have one of the six :
		// 1. conversation_id in the body or in the URL
		// 2. message_id in the body or in the URL
		// 3. comment_id in the body or in the URL

		// First, we check for the conversation_id case :

		// Check if the conversation_id is in the URL
		SpecifiedConversationID = ps.ByName("conversationId")
		if SpecifiedConversationID == "" {
			rt.baseLogger.Infoln("No conversation_id in the URL. Checking the request body...")
			// Check if the conversation_id is in the request body
			var ConversationBody struct {
				ConversationID string `json:"conversationId"`
			}
			err = json.NewDecoder(r.Body).Decode(&ConversationBody)
			if err != nil {
				rt.baseLogger.Infoln("No conversation_id in the request body. proceeding with message_id check...")
			} else {
				SpecifiedConversationID = ConversationBody.ConversationID
			}
			// Reset the request body so it can be read again
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// if the conversation_id is still empty, we proceed with the message_id check.
		if SpecifiedConversationID == "" {
			messageID := ps.ByName("messageId")
			if messageID == "" {
				// if the message_id is still empty, we check if the message_id is in the request body.
				rt.baseLogger.Infoln("No message_id in the URL. Checking the request body...")
				var MessageBody struct {
					MessageID string `json:"messageId"`
				}
				err = json.NewDecoder(r.Body).Decode(&MessageBody)
				if err != nil {
					rt.baseLogger.Infoln("No message_id in the request body either... proceeding with comment_id check...")
				} else {
					messageID = MessageBody.MessageID
				}
				// Reset the request body so it can be read again
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}
			if messageID != "" {
				// if the message_id is not empty, we get the conversation_id from the message_id.
				// we will use rt.db.GetMessageInfo(messageID, token)
				// this will return the conversation_id, sender_id, content, timestamp, and the message_id.
				// we only need the conversation_id.
				SpecifiedConversationID, _, _, _, _, _, err = rt.db.GetMessageInfo(messageID, token)
				if err != nil {
					http.Error(w, "Failed to get message info", http.StatusInternalServerError)
					return
				}
			}
			// if the message_id is still empty, we proceed with the comment_id check.
		}

		if SpecifiedConversationID == "" {
			commentID := ps.ByName("commentId")
			if commentID == "" {
				rt.baseLogger.Infoln("No comment_id in the URL. Checking the request body...")
				// if the comment_id is still empty, we check if the comment_id is in the request body.
				var CommentBody struct {
					CommentID string `json:"commentId"`
				}
				err = json.NewDecoder(r.Body).Decode(&CommentBody)
				if err != nil {
					rt.baseLogger.Infoln("No comment_id in the request body either...")
				} else {
					commentID = CommentBody.CommentID
				}
				// Reset the request body so it can be read again
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}
			if commentID != "" {
				var CommentMessageID string
				CommentMessageID, _, _, _, err = rt.db.GetCommentInfo(commentID)
				if err != nil {
					http.Error(w, "Failed to get comment info", http.StatusInternalServerError)
					return
				}
				SpecifiedConversationID, _, _, _, _, _, err = rt.db.GetMessageInfo(CommentMessageID, token)
				if err != nil {
					http.Error(w, "Failed to get message info", http.StatusInternalServerError)
					return

				}
			}
		}

		// if the conversation_id is still empty, we return an error.
		if SpecifiedConversationID == "" {
			rt.baseLogger.Infoln("No conversation_id in the request body or URL.")
			http.Error(w, "No conversation_id could be found from the given request.", http.StatusBadRequest)
			return
		}

		// Check if the user is a participant in the conversation
		participants, err := rt.db.GetParticipants(SpecifiedConversationID)

		isParticipant := false
		for _, participant := range participants {
			if participant == token {
				isParticipant = true
				break
			}
		}

		if err != nil {
			http.Error(w, "Failed to check user in conversation", http.StatusInternalServerError)
			return
		}
		if !isParticipant {
			rt.baseLogger.Infoln("User " + token + " is NOT a participant of conversation " + SpecifiedConversationID + ".")
			http.Error(w, "User "+token+" is not a participant of conversation "+SpecifiedConversationID, http.StatusForbidden)
			return
		}

		rt.baseLogger.Infoln("User " + token + " is a participant of conversation " + SpecifiedConversationID + ". Operation allowed.")

		// Call the next handler if the user is a participant
		next(w, r, ps)
	}
}
