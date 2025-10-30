package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type SetMyUserNameRequest struct {
	New_name string `json:"new_name"`
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var req SetMyUserNameRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON:", err.Error())
		return
	}
	rt.baseLogger.Info("received request to change user name to '" + req.New_name + "'")

	// Verify the db is ok
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check the new name pattern
	if len(req.New_name) < 3 || len(req.New_name) > 16 {
		rt.baseLogger.Error("new user name '" + req.New_name + "' is too long or too short (min 3, max 16)")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("new user name '" + req.New_name + "' is too long or too short (min 3, max 16)")
		return
	}

	// Check whether the new name already exists!
	id, _ := rt.db.GetUserId(req.New_name)
	if id != 0 {
		rt.baseLogger.Error("the user name '" + req.New_name + "' already exists")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("the user name '" + req.New_name + "' already exists")
		return
	}

	// Update the name
	err = rt.db.SetMyUserName(user.Id, req.New_name)
	if err != nil {
		rt.baseLogger.Error("error trying to update the user name (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("name succesfully updated.")
}
