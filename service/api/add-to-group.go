package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AddToGroupRequest struct {
	User_id_to_add   int
	User_name_to_add string `json:"user_name_to_add"`
	Group_id         int
	Group_name       string `json:"group_name"`
	Conversation_id  int
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var req AddToGroupRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON (", err.Error()+")")
		return
	}
	rt.baseLogger.Info("received request to add user ", req.User_name_to_add+" to group "+req.Group_name)

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the user id to add
	req.User_id_to_add, _ = rt.db.GetUserId(req.User_name_to_add)
	if req.User_id_to_add == 0 {
		rt.baseLogger.Error("the user `", req.User_name_to_add, "` does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the group id based on the group name
	req.Group_id, _, req.Conversation_id, _ = rt.db.GetGroupDataFromName(req.Group_name)
	if req.Group_id == 0 {
		// Create a new conversation
		req.Conversation_id, _ = rt.db.GetNextConversationId()
		err = rt.db.CreateConversation(req.Conversation_id, user.Id, req.User_id_to_add)
		if err != nil {
			rt.baseLogger.Error("error while creating a new conversation for a group (", err.Error()+")")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create a new group and conversation
		req.Group_id, _ = rt.db.GetNextGroupId()
		err = rt.db.CreateGroup(req.Group_id, req.Group_name, ``, req.Conversation_id)
		if err != nil {
			rt.baseLogger.Error("error while creating a new conversation for a group (", err.Error()+")")
			rt.baseLogger.Error("the table groups and conversations could be in an anomalous state")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rt.baseLogger.Info("the group `", req.Group_name, "` does not exist so a new one has been created with id "+fmt.Sprint(req.Group_id)+" and conversation id ", fmt.Sprint(req.Conversation_id))
	} else {
		// Add user to existing group, that is its conversation
		err = rt.db.AddUserToConversation(req.Conversation_id, req.User_id_to_add)
		if err != nil {
			_ = json.NewEncoder(w).Encode("l'utente è stato già aggiunto al gruppo")
			// rt.baseLogger.Error("error while adding user to group (", err.Error()+")")
			// w.WriteHeader(http.StatusInternalServerError)
			// return
		}
	}
	w.WriteHeader(http.StatusCreated)

	// Exit
	rt.baseLogger.Info("user added to group")
}
