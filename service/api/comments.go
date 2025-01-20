package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// add a comment to a message
func (rt *_router) AddComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// message_id is passed in the request body
	var comment struct {
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	MessageID := ps.ByName("messageId")

	// Extract the token from the Authorization header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	rt.baseLogger.Infof("Received token : %v", token)
	// further token checks are done in the wrapper function so we don't need to check it here

	// add the comment to the database
	err = rt.db.AddComment(MessageID, token, comment.Content)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error adding the comment to the database")
		http.Error(w, "Failed to add the comment", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Infof("Comment added to message %v", MessageID)
	w.WriteHeader(http.StatusCreated)
}

// get all comments for a message
func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// message_id is passed in the json payload

	MessageID := ps.ByName("messageId")

	// get the comments from the database
	commentIDs, _, senderIDs, contents, timestamps, err := rt.db.GetComments(MessageID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the comments from the database")
		http.Error(w, "Failed to get the comments", http.StatusInternalServerError)
		return
	}

	// create the response
	type commentResponse struct {
		CommentID string `json:"commentId"`
		SenderID  string `json:"senderId"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}
	comments := make([]commentResponse, 0)
	for i := range commentIDs {
		comments = append(comments, commentResponse{
			CommentID: commentIDs[i],
			SenderID:  senderIDs[i],
			Content:   contents[i],
			Timestamp: timestamps[i],
		})
	}

	// write the response
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding the response")
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// delete a comment from a message
func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// comment_id is passed in the json payload

	CommentID := ps.ByName("commentId")

	// delete the comment from the database
	err := rt.db.DeleteComment(CommentID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error deleting the comment from the database")
		http.Error(w, "Failed to delete the comment", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Infof("Comment %v deleted", CommentID)
	w.WriteHeader(http.StatusOK)
}

// get information about a comment
func (rt *_router) GetCommentInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// comment_id is passed in the json payload

	CommentID := ps.ByName("commentId")

	// get the comment info from the database
	messageID, senderID, content, timestamp, err := rt.db.GetCommentInfo(CommentID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the comment info from the database")
		http.Error(w, "Failed to get the comment info", http.StatusInternalServerError)
		return
	}

	sender_username, err := rt.db.GetUserName(senderID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the sender's username")
		http.Error(w, "Failed to get the sender's username", http.StatusInternalServerError)
		return
	}

	// create the response
	type commentResponse struct {
		MessageID string `json:"messageId"`
		Sender    string `json:"sender"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}
	resp := commentResponse{
		MessageID: messageID,
		Sender:    sender_username,
		Content:   content,
		Timestamp: timestamp,
	}

	// write the response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding the response")
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
}
