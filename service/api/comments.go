package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// add a comment to a message
func (rt *_router) AddComment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// message_id is passed in the request body
	var comment struct {
		MessageID string `json:"message_id"`
		Content   string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Extract the token from the Authorization header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	rt.baseLogger.Infof("Received token : %v", token)
	// further token checks are done in the wrapper function so we don't need to check it here

	// add the comment to the database
	err = rt.db.AddComment(comment.MessageID, token, comment.Content)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error adding the comment to the database")
		http.Error(w, "Failed to add the comment", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Infof("Comment added to message %v", comment.MessageID)
	w.WriteHeader(http.StatusOK)
}

// get all comments for a message
func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// message_id is passed in the json payload
	var message struct {
		MessageID string `json:"message_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// get the comments from the database
	commentIDs, _, senderIDs, contents, timestamps, err := rt.db.GetComments(message.MessageID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the comments from the database")
		http.Error(w, "Failed to get the comments", http.StatusInternalServerError)
		return
	}

	// create the response
	type commentResponse struct {
		CommentID string `json:"comment_id"`
		SenderID  string `json:"sender_id"`
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

	rt.baseLogger.Infof("Comments for message %v sent", message.MessageID)
	w.WriteHeader(http.StatusOK)
}

// delete a comment from a message
func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// comment_id is passed in the json payload
	var comment struct {
		CommentID string `json:"comment_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// delete the comment from the database
	err = rt.db.DeleteComment(comment.CommentID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error deleting the comment from the database")
		http.Error(w, "Failed to delete the comment", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Infof("Comment %v deleted", comment.CommentID)
	w.WriteHeader(http.StatusOK)
}

// get information about a comment
func (rt *_router) GetCommentInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// comment_id is passed in the json payload
	var comment struct {
		CommentID string `json:"comment_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// get the comment info from the database
	messageID, senderID, content, timestamp, err := rt.db.GetCommentInfo(comment.CommentID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the comment info from the database")
		http.Error(w, "Failed to get the comment info", http.StatusInternalServerError)
		return
	}

	// create the response
	type commentResponse struct {
		MessageID string `json:"message_id"`
		SenderID  string `json:"sender_id"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}
	resp := commentResponse{
		MessageID: messageID,
		SenderID:  senderID,
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

	rt.baseLogger.Infof("Comment info for comment %v sent", comment.CommentID)
}
