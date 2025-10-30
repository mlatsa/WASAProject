package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Contains info from the request body AND from local variable
type SetMyPhotoRequest struct {
	User_id int
	Photo   string `json:"photo"`
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Authenticate user through session token
	var user User
	var err error
	user, err = AuthenticateUser(r)
	if err != nil {
		rt.baseLogger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	rt.baseLogger.Info("authenticated user `", user.Name, "`, id ", user.Id, "")

	// Read the request
	var photoreq SetMyPhotoRequest
	err = json.NewDecoder(r.Body).Decode(&photoreq)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON:", err.Error())
		return
	}
	rt.baseLogger.Info("received a new profile photo for the user `", user.Name, "`")

	// Verify the db is ok
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the photo
	err = rt.db.SetMyPhoto(user.Id, photoreq.Photo)
	if err != nil {
		rt.baseLogger.Error("error trying to update the user photo (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("Photo succesfully updated.")
}
