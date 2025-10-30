package api

import (
	"WasaText/service/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type GetMembersRequest struct {
	Group_id int
}

func (rt *_router) getMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var req GetMembersRequest
	req.Group_id, err = strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/groups/"))
	if err != nil {
		rt.baseLogger.Error("error while reading group id from the path (" + err.Error() + ")")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("received request to get members of the group with id ", req.Group_id)

	// Get the members
	var users []database.User
	users, err = rt.db.GetMembers(req.Group_id)
	if err != nil {
		rt.baseLogger.Error("error while retriving the members of the group with id ", req.Group_id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("retrieved ", len(users), " members")

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	err_encoder := json.NewEncoder(w).Encode(users)
	if err_encoder != nil {
		rt.baseLogger.Error("error during the encoding of the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
