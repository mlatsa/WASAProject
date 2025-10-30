package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Contains info from the request body AND from the path AND from local variable
type ForwardMessageRequest struct {
	From_msg_id                int
	Msg                        string
	Conversation_id            int
	To_user_name_or_group_name string `json:"to_user_name_or_group_name"`
	Is_group                   bool   `json:"is_group"`
	To_user_id_or_group_id     int
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Authenticate user through session token
	var err error
	var user User
	user, err = AuthenticateUser(r)
	if err != nil {
		rt.baseLogger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	rt.baseLogger.Info("authenticated user `", user.Name, "`, id ", user.Id, "")

	// Retrieve the message_id
	var req ForwardMessageRequest
	// TO DO: should verify that the message belongs to a user's conversation
	req.From_msg_id, err = strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/messages/"))
	if err != nil {
		rt.baseLogger.Error("error while reading message id from the path (" + err.Error() + ")")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the request body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON (", err.Error()+")")
		return
	}
	rt.baseLogger.Info("received request to forward message with id ", req.From_msg_id, " to the user or group ", req.To_user_name_or_group_name)

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retrieve the message
	req.Msg, err = rt.db.GetMessageText(req.From_msg_id)
	if err != nil {
		rt.baseLogger.Error("error retrieving message id ", req.From_msg_id, " ("+err.Error()+")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Manage the user case
	if !req.Is_group {
		// Get the receiver id
		req.To_user_id_or_group_id, _ = rt.db.GetUserId(req.To_user_name_or_group_name)
		if req.To_user_id_or_group_id == 0 {
			rt.baseLogger.Error("the receipent of the message `", req.To_user_name_or_group_name, "` does not exist ")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Retrieve the conversation_id based on the user id's, if 0 was provided
		if req.Conversation_id == 0 {
			req.Conversation_id, _ = rt.db.GetConversationId(user.Id, req.To_user_id_or_group_id)
		}

		// If the conversation does not exist, create it
		if req.Conversation_id == 0 {
			req.Conversation_id, err = rt.db.GetNextConversationId()
			if err != nil {
				rt.baseLogger.Error("error while retrieving the next conversation id (", err.Error()+")")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = rt.db.CreateConversation(req.Conversation_id, user.Id, req.To_user_id_or_group_id)
			if err != nil {
				rt.baseLogger.Error("error while creating the new conversation (", err.Error()+")")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			rt.baseLogger.Info("the conversation does not exist so a new one is created with conversation id ", req.Conversation_id)
		} else {
			rt.baseLogger.Info("retrieved conversation id ", req.Conversation_id)
		}
	} else {
		// Manage the group case

		// Get the receiver id
		req.To_user_id_or_group_id, _, req.Conversation_id, _ = rt.db.GetGroupDataFromName(req.To_user_name_or_group_name)
		if req.To_user_id_or_group_id == 0 {
			rt.baseLogger.Error("the receipent of the message `", req.To_user_name_or_group_name, "` does not exist ")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Forward the message
	err = rt.db.SendMessage(req.Conversation_id, user.Id, req.Msg, req.From_msg_id)
	if err != nil {
		rt.baseLogger.Error("error while forwarding a message (" + err.Error() + ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	// Exit
	rt.baseLogger.Info("message forwarded")
}
