package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) changeUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	rt.baseLogger.Infof("Received token : %v", token)

	// new username is passed in the request body
	var newUsername struct {
		Username string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// add the new username to the database, with the same user_id as the old username.
	err = rt.db.UpdateUser(newUsername.Username, token, "")
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error updating the user in the database")
		http.Error(w, "Failed to update the user", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Infof("User %v updated to %v", token, newUsername.Username)
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) changeUserPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	rt.baseLogger.Infof("Received token : %v", token)

	// new photo is passed in the request body
	var newPhoto struct {
		Photo string `json:"photo"`
	}
	err := json.NewDecoder(r.Body).Decode(&newPhoto)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// add the new photo to the database, with the same user_id as the old photo.
	err = rt.db.UpdateUser("", token, newPhoto.Photo)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error updating the user in the database")
		http.Error(w, "Failed to update the user", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Infof("User %v updated with photo %v", token, newPhoto.Photo)
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	rt.baseLogger.Infof("Received token : %v", token)

	// get all users from the database
	users, err := rt.db.GetAllUsers()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the users from the database")
		http.Error(w, "Failed to get the users", http.StatusInternalServerError)
		return
	}

	// return the users as a json payload
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding the users")
		http.Error(w, "Failed to encode the users", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) checkUserExistence(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	rt.baseLogger.Infof("Received token : %v", token)

	// user_id is passed in the URL
	userID := ps.ByName("userId")

	// check if the user exists in the database
	exists, err := rt.db.CheckUserExistence(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error checking the user existence in the database")
		http.Error(w, "Failed to check the user existence", http.StatusInternalServerError)
		return
	}

	// return the result as a json payload
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Exists bool `json:"exists"`
	}{
		Exists: exists,
	})
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding the existence")
		http.Error(w, "Failed to encode the existence", http.StatusInternalServerError)
		return
	}
}

// get a specified user's profile details from the database (username, photo)
func (rt *_router) getUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// from the path
	username := ps.ByName("username")

	// get the user's profile details from the database
	photo, err := rt.db.GetUserPhoto(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting the user's photo from the database")
		http.Error(w, "Failed to get the user's photo", http.StatusInternalServerError)
		return
	}

	// return the user's profile details as a json payload
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Photo string `json:"photo"`
	}{
		Photo: photo,
	})
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding the user's profile")
		http.Error(w, "Failed to encode the user's profile", http.StatusInternalServerError)
		return
	}
}
