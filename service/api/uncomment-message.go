package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retrieve the message id
	var msg_id int
	msg_id, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/messages/:"), "/comments"))
	rt.baseLogger.Info("received request to delete comment on message ", msg_id)

	// Delete the comment
	err = rt.db.UncommentMessage(msg_id, user.Id)
	if err != nil {
		rt.baseLogger.Error("error while uncommenting the message (", err.Error()+")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Exit
	rt.baseLogger.Info("comment deleted")
}
