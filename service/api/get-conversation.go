package api

import (
	"WasaText/service/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Get conversation AND mark messages as seen by the caller user
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Retrieve the conversation_id and check if it belongs to the user
	conversation_id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/conversations/"))
	if err != nil {
		rt.baseLogger.Error("error reading conversation_id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("received request to get conversation with id ", conversation_id)
	belongTo, _ := rt.db.IsUserInConversation(user.Id, conversation_id)
	if !belongTo {
		rt.baseLogger.Error("the conversation_id do not belong to the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the conversation
	var messages []database.Message
	messages, err = rt.db.GetConversation(conversation_id)
	if err != nil {
		rt.baseLogger.Error("error trying to get conversations (", err.Error(), ")")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("retrieved ", len(messages), " messages")

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		rt.baseLogger.Error("error during the encoding of the json (", err.Error(), ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Mark the all the messages of the conversatin as seen by the user!
	err = rt.db.MarkMessagesSeen(user.Id, conversation_id)
	if err != nil {
		rt.baseLogger.Error("issue while marking messages seen (", err.Error(), ")")
	}
}
