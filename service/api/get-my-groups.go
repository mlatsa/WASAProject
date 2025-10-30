package api

import (
	"WasaText/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func (rt *_router) getMyGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Get the groups that belongs to the user
	var groups []database.Group
	groups, err = rt.db.GetMyGroups(user.Id)
	if err != nil {
		rt.baseLogger.Error("error trying to get group list (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("retrieved ", len(groups), " groups")

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		rt.baseLogger.Error("error during the encoding of the json (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
