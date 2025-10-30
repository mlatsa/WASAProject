package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	msg_id, _ = strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/messages/"))
	rt.baseLogger.Info("received request to delete message ", msg_id)

	// Delete the message
	var result sql.Result
	result, err = rt.db.DeleteMessage(msg_id)
	if err != nil {
		rt.baseLogger.Error("error while deleting the message ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var rowsAffected int64
	rowsAffected, _ = result.RowsAffected()
	if rowsAffected == 0 {
		rt.baseLogger.Error("the message id ", msg_id, " does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Exit
	rt.baseLogger.Info("message deleted")
}
