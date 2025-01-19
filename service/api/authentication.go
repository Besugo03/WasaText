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

// checkUserInConversation is a middleware that checks if the user is a participant in the conversation, and if the conversation exists, before calling the next handler.
func (rt *_router) checkUserInConversation(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extract the token from the Authorization header
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		rt.baseLogger.Infof("Received token : %v", token)

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

		// Check the JSON payload for the conversation_id
		type sendMessageRequest struct {
			ConversationID string `json:"conversation_id"`
		}
		var req sendMessageRequest
		if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&req); err != nil {
			// the requestID may be passed as a URL parameter, so try to read it from there
			req.ConversationID = ps.ByName("conversation_id")
			if req.ConversationID == "" {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}
		}

		specifiedConversationID := req.ConversationID

		// Check if the conversation exists
		_, _, _, _, err = rt.db.GetConversationInfo(specifiedConversationID, token)
		if err != nil {
			// it might be that the conversation is not visible in the request.
			// this might happen if it's a message-specific or comment-specific request.

			// catching the message edgecase
			type messageRequest struct {
				MessageID string `json:"message_id"`
			}
			var req messageRequest
			if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&req); err != nil {
				http.Error(w, "Invalid request payload. Did you check if the conversation exists, or if the user is part of it?", http.StatusBadRequest)
				return
			}

			// get the conversation_id from the message_id
			conversationID, _, _, _, _, _, err := rt.db.GetMessageInfo(req.MessageID, token)
			if conversationID == "" {
				// if we STILL didnt get any info on what the conversation is, it's probably a comment request.
				type commentRequest struct {
					CommentID string `json:"comment_id"`
				}
				var req commentRequest
				if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&req); err != nil {
					http.Error(w, "Invalid request payload. Did you check if the conversation exists, or if the user is part of it?", http.StatusBadRequest)
					return
				}
				// get the conversation_id from the comment_id
				messageID, _, _, _, err := rt.db.GetCommentInfo(req.CommentID)
				if err != nil {
					rt.baseLogger.WithError(err).Errorf("Error querying the database for comment %v", req.CommentID)
					http.Error(w, "Comment not found", http.StatusNotFound)
					return
				}
				conversationID, _, _, _, _, _, err = rt.db.GetMessageInfo(messageID, token)
				rt.baseLogger.Infoln("Conversation ID from message ID " + messageID + " is " + conversationID)

				if err != nil {
					rt.baseLogger.WithError(err).Errorf("Error querying the database for message %v", messageID)
					http.Error(w, "Message not found", http.StatusNotFound)
					return
				}
			}
			if err != nil && conversationID == "" {
				rt.baseLogger.WithError(err).Errorf("Error querying the database for message %v", req.MessageID)
				http.Error(w, "Message not found", http.StatusNotFound)
				return
			}
			specifiedConversationID = conversationID
			rt.baseLogger.Infof("Conversation ID from message ID %v is %v", req.MessageID, conversationID)

			if specifiedConversationID == "" {
				http.Error(w, "Conversation not found", http.StatusNotFound)
				return
			}
		}

		// Check if the user is a participant in the conversation
		participants, err := rt.db.GetParticipants(specifiedConversationID)

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
			rt.baseLogger.Infoln("User " + token + " is NOT a participant of conversation " + specifiedConversationID + ".")
			http.Error(w, "User "+token+" is not a participant of conversation "+specifiedConversationID, http.StatusForbidden)
			return
		}

		rt.baseLogger.Infoln("User " + token + " is a participant of conversation " + specifiedConversationID + ". Operation allowed.")

		// Call the next handler if the user is a participant
		next(w, r, ps)
	}
}
