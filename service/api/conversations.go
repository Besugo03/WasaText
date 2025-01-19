package api

import (
	"encoding/json"
	"fmt"
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
	rt.baseLogger.Infof("Getting conversation info for conversation %v", ps.ByName("conversation_id"))
	fmt.Println("Getting conversation info for conversation", ps.ByName("conversation_id"))
	conversationID := ps.ByName("conversation_id")
	conversationID, isGroup, name, picture, err := rt.db.GetConversationInfo(conversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting conversation info for conversation %v", conversationID)
		http.Error(w, "Failed to get conversation info", http.StatusInternalServerError)
		return
	}

	participants, err := rt.db.GetParticipants(conversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting participants for conversation %v", conversationID)
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"conversation_id": conversationID,
		"is_group":        isGroup,
		"name":            name,
		"picture":         picture,
		"participants":    participants,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// in the body we will have:
	// conversation_id, content, is_forwarded
	type sendMessageRequest struct {
		ConversationID string `json:"conversation_id"`
		Content        string `json:"content"`
		IsForwarded    bool   `json:"is_forwarded"`
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.baseLogger.Infoln(req)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("Authorization")
	userID = userID[7:] // strip the "Bearer " prefix

	// first, we must check if the conversation exists
	_, _, _, _, err := rt.db.GetConversationInfo(req.ConversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting conversation info for conversation %v", req.ConversationID)
		http.Error(w, "Conversation does not exist", http.StatusNotFound)
		return
	}
	// then, we need to check if the user is a participant in the conversation
	participants, err := rt.db.GetParticipants(req.ConversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting participants for conversation %v", req.ConversationID)
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}
	isParticipant := false
	for _, participant := range participants {
		if participant == userID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		http.Error(w, "User is not a participant in the conversation", http.StatusForbidden)
		return
	}

	messageID, err := rt.db.SendMessage(req.ConversationID, userID, req.Content, req.IsForwarded)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error sending message to conversation %v", req.ConversationID)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message_id": messageID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) CreateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// in the header, we will get the participants to the conversation.
	// is_group is automatic (true if more than 2 participants)
	// name and picture are expected to be in the request body only if is_group is true
	type createConversationRequest struct {
		Name         string   `json:"name"`
		Picture      string   `json:"picture"`
		Participants []string `json:"participants"`
	}

	var req createConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// since the participants are passed as usernames, we need to convert them to user IDs
	participant_ids := make([]string, len(req.Participants))
	for i, participant := range req.Participants {
		userID, err := rt.db.GetUserID(participant)
		if err != nil {
			rt.baseLogger.WithError(err).Errorf("error getting user ID for user %v", participant)
			http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
			return
		}
		participant_ids[i] = userID
	}

	// lastly, we add the user who created the conversation
	userID := r.Header.Get("Authorization")
	userID = userID[7:] // strip the "Bearer " prefix
	participant_ids = append(participant_ids, userID)

	isGroup := len(participant_ids) > 2
	conversationID, err := rt.db.CreateConversation(isGroup, req.Name, req.Picture)
	if err != nil {
		rt.baseLogger.WithError(err).Error("error creating conversation")
		http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	for _, participant := range participant_ids {
		err := rt.db.AddParticipant(conversationID, participant)
		if err != nil {
			rt.baseLogger.WithError(err).Errorf("error adding participant %v to conversation %v", participant, conversationID)
			http.Error(w, "Failed to add participant", http.StatusInternalServerError)
			return
		}
	}

	resp := map[string]interface{}{
		"conversation_id": conversationID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) UpdateConversationName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// in the body we will have:
	// name
	// conversation_id
	type updateConversationNameRequest struct {
		Name           string `json:"name"`
		ConversationID string `json:"conversation_id"`
	}

	var req updateConversationNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	conversationID := req.ConversationID

	err := rt.db.UpdateConversation(conversationID, false, req.Name, "")
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error updating conversation %v", conversationID)
		http.Error(w, "Failed to update conversation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) UpdateConversationPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// in the body we will have:
	// photo
	// conversation_id
	type updateConversationPhotoRequest struct {
		Photo          string `json:"photo"`
		ConversationID string `json:"conversation_id"`
	}

	var req updateConversationPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := rt.db.UpdateConversation(req.ConversationID, false, "", req.Photo)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error updating conversation %v", req.ConversationID)
		http.Error(w, "Failed to update conversation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (rt *_router) AddParticipant(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// conversation ID and the user ID are expected to be in the request body
	type addParticipantRequest struct {
		ConversationID string `json:"conversation_id"`
		Username       string `json:"username"`
	}
	var req addParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	conversationID := req.ConversationID

	userID, err := rt.db.GetUserID(req.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting user ID for user %v", req.Username)
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	// first, we must check if the conversation exists
	_, _, _, _, err = rt.db.GetConversationInfo(conversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting conversation info for conversation %v", conversationID)
		http.Error(w, "Conversation does not exist", http.StatusNotFound)
		return
	}

	// then, we need to check if the user is a participant in the conversation already
	participants, err := rt.db.GetParticipants(conversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting participants for conversation %v", conversationID)
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}
	isParticipant := false
	for _, participant := range participants {
		if participant == userID {
			isParticipant = true
			break
		}
	}
	if isParticipant {
		http.Error(w, "User is already a participant in the conversation", http.StatusConflict)
		return
	}

	err = rt.db.AddParticipant(conversationID, userID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error adding participant %v to conversation %v", userID, conversationID)
		http.Error(w, "Failed to add participant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rt *_router) RemoveParticipant(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// conversation ID and the user ID are expected to be in the request body
	// all necessary checks are done from the checkUserInConversation middleware.
	type removeParticipantRequest struct {
		ConversationID string `json:"conversation_id"`
		Username       string `json:"username"`
	}
	var req removeParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	conversationID := req.ConversationID

	userID, err := rt.db.GetUserID(req.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting user ID for user %v", req.Username)
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	// check if the user we want to remove is part of the conversation
	participants, err := rt.db.GetParticipants(conversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting participants for conversation %v", conversationID)
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}
	isParticipant := false
	for _, participant := range participants {
		if participant == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		http.Error(w, "User to remove is not a participant in the conversation", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveParticipant(conversationID, userID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error removing participant %v from conversation %v", userID, conversationID)
		http.Error(w, "Failed to remove participant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// this basically just gets a message_id from the database, and sends it in another conversation, with the is_forwarded flag set to true.
func (rt *_router) ForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// in the body we will have:
	// conversation_id, message_id
	type forwardMessageRequest struct {
		ConversationID string `json:"conversation_id"`
		MessageID      string `json:"message_id"`
	}

	var req forwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("Authorization")
	userID = userID[7:] // strip the "Bearer " prefix

	// get all the messages
	messageIDs, _, _, contents, _, err := rt.db.GetMessages(req.ConversationID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting messages for conversation %v", req.ConversationID)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	// check if the message exists
	messageExists := false
	for _, messageID := range messageIDs {
		if messageID == req.MessageID {
			messageExists = true
			break
		}
	}
	if !messageExists {
		http.Error(w, "Message does not exist", http.StatusNotFound)
		return
	}

	// for the specified message, get the content
	var content string = ""
	for i, messageID := range messageIDs {
		if messageID == req.MessageID {
			content = contents[i]
			break
		}
	}

	// send the message to the conversation
	messageID, err := rt.db.SendMessage(req.ConversationID, userID, content, true)

	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error sending message to conversation %v", req.ConversationID)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message_id": messageID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// in the body we will have:
	// conversation_id, message_id
	type deleteMessageRequest struct {
		MessageID string `json:"message_id"`
	}

	var req deleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// check if the message exists
	_, _, _, _, err := rt.db.GetMessageInfo(req.MessageID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error getting message info for message %v", req.MessageID)
		http.Error(w, "Failed to get message info", http.StatusInternalServerError)
		return
	}

	// delete the message
	err = rt.db.DeleteMessage(req.MessageID)
	if err != nil {
		rt.baseLogger.WithError(err).Errorf("error deleting message %v", req.MessageID)
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
