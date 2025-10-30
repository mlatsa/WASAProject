package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Contains info from the request body AND from local variable
type SetGroupPhotoRequest struct {
	Group_id int
	Photo    string `json:"photo"`
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var req SetGroupPhotoRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON:", err.Error())
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	rt.baseLogger.Info("received a new group photo for the group id `", req.Group_id, "`")

	// Verify the db is ok
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the group photo
	err = rt.db.SetGroupPhoto(req.Group_id, req.Photo)
	if err != nil {
		rt.baseLogger.Error("error trying to update the group photo (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("Group photo succesfully updated.")
}
