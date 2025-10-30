package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

var sessions = make(map[int64]User)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read the request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("Error decoding JSON:", err)
		return
	}
	rt.baseLogger.Info("received login request for user '" + user.Name + "'")

	// Verify the db is ok
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check the user name pattern
	if len(user.Name) < 3 || len(user.Name) > 16 {
		rt.baseLogger.Error("user name '" + user.Name + "' is too long or too short.")
		w.WriteHeader(http.StatusBadRequest)
		// w.Write() // w.Write must go after WriteHeader otherwise StatusOK is written!
		_ = json.NewEncoder(w).Encode("user name '" + user.Name + "' is too long or too short.")
		return
	}

	// Create a new user if necessary
	var new_user bool = false
	user.Id, user.Photo, err = rt.db.GetUserIdAndPhoto(user.Name)
	if err != nil {
		rt.baseLogger.Error("error while retrieving user data " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.Id == 0 {
		new_user = true
		user.Id, _ = rt.db.GetNextUserId()
		rt.baseLogger.Info("new user! assigned id ", user.Id)
		err = rt.db.CreateUser(user.Id, user.Name)
		if err != nil {
			rt.baseLogger.Error("error while creating a new user " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Generate a random token
	var session_token int64 = rand.Int63()
	rt.baseLogger.Info("session token ", session_token)
	sessions[session_token] = user

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+fmt.Sprint(session_token))
	if new_user {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Send the user data as content
	err_encoder := json.NewEncoder(w).Encode(user)
	if err_encoder != nil {
		rt.baseLogger.Error("error during the encoding of the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Exit
	rt.baseLogger.Info("login done successfully")
}
