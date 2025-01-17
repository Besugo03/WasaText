package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Header.Get("Authorization")
	userID = userID[7:] // strip the "Bearer " prefix

	// Retrieve conversations for the user (this is just a placeholder, replace with actual logic)
	// conversations := []string{"Conversation 1", "Conversation 2"}
	conversations, err := rt.db.GetUserConversations(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting conversations for user %v", userID)
		http.Error(w, "Failed to get conversations", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"user_id":       userID,
		"conversations": conversations,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) GetConversationInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conversationID := ps.ByName("conversation_id")
	conversationID, isGroup, name, picture, err := rt.db.GetConversationInfo(conversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting conversation info for conversation %v", conversationID)
		http.Error(w, "Failed to get conversation info", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"conversation_id": conversationID,
		"is_group":        isGroup,
		"name":            name,
		"picture":         picture,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// conversationID and the message content are expected to be in the request body
	type messageRequest struct {
		ConversationID string `json:"conversation_id"`
		Message        string `json:"message"`
	}
	var req messageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Send the message to the conversation
	// TODO implement this.
}

func (rt *_router) AddParticipant(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// conversation ID and the user ID are expected to be in the request body
	type addParticipantRequest struct {
		ConversationID string `json:"conversation_id"`
		UserID         string `json:"user_id"`
	}
	var req addParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	conversationID := req.ConversationID
	userID := req.UserID

	err := rt.db.AddParticipant(conversationID, userID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error adding participant %v to conversation %v", userID, conversationID)
		http.Error(w, "Failed to add participant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
