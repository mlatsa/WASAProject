package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type LeaveGroupRequest struct {
	Group_id int
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Retrieve the group_id
	var req LeaveGroupRequest
	req.Group_id, err = strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/groups/"))
	if err != nil {
		rt.baseLogger.Error("error while reading group id from the path (" + err.Error() + ")")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("received request leave the group with id ", req.Group_id)

	// Leave the group
	err = rt.db.LeaveGroup(user.Id, req.Group_id)
	if err != nil {
		rt.baseLogger.Error("error while leaving the group with id `", req.Group_id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Exit
	rt.baseLogger.Info("Left group")
}
