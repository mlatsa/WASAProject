package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CommentMessageRequest struct {
	Msg_id   int    `json:"msg_id"`
	Reaction string `json:"reaction"`
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var req CommentMessageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON (", err.Error()+")")
		return
	}
	rt.baseLogger.Info("received request to comment message ", req.Msg_id)

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TO DO: should verify the message belongs to a user’s conversation

	// Before commenting remove the existing user comment
	err = rt.db.UncommentMessage(req.Msg_id, user.Id)
	if err != nil {
		rt.baseLogger.Error("error while uncommenting a message (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Comment the message
	_, err = rt.db.CommentMessage(req.Msg_id, user.Id, req.Reaction)
	if err != nil {
		rt.baseLogger.Error("error while commenting a message (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	// Exit
	rt.baseLogger.Info("message commented")
}
