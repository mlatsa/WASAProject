package api

import (
	"WasaText/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Get the list of users
	var users []database.User
	users, err = rt.db.GetAllUsers()
	if err != nil {
		rt.baseLogger.Error("error trying to get the list of users (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("retrieved ", len(users), " users")

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	err_encoder := json.NewEncoder(w).Encode(users)
	if err_encoder != nil {
		rt.baseLogger.Error("error during the encoding of the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
