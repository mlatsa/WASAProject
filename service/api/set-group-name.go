package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type SetGroupNameRequest struct {
	Group_id       int
	New_group_name string `json:"new_group_name"`
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Read the group_id
	var req SetGroupNameRequest
	req.Group_id, err = strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/groups/"))
	if err != nil {
		rt.baseLogger.Error("error reading message id from the path (" + err.Error() + ")")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("error reading message id from the path (" + err.Error() + ")")
		return
	}

	// Read the request body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON:", err.Error())
		return
	}
	rt.baseLogger.Info("received request to change group name for group id '", req.Group_id, "' to '"+req.New_group_name+"'")

	// Verify the db is ok
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check the new name pattern
	if len(req.New_group_name) < 3 || len(req.New_group_name) > 16 {
		rt.baseLogger.Error("new user name '" + req.New_group_name + "' is too long or too short (min 3, max 16)")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("new user name '" + req.New_group_name + "' is too long or too short (min 3, max 16)")
		return
	}

	// Check whether the new name already exists!
	id, _, _, _ := rt.db.GetGroupDataFromName(req.New_group_name)
	if id != 0 {
		rt.baseLogger.Error("the group name '" + req.New_group_name + "' already exists")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("the group name '" + req.New_group_name + "' already exists")
		return
	}

	// Update the name
	err = rt.db.SetGroupName(req.Group_id, req.New_group_name)
	if err != nil {
		rt.baseLogger.Error("error trying to update the group name (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("name succesfully updated.")
}
