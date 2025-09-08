package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// ----- health -----
func (rt *Router) health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ----- user profile -----
type putUsernameBody struct {
	Username string `json:"username"`
}

func (rt *Router) putUserUsername(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	var b putUsernameBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil || len(b.Username) < 3 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"detail": "invalid username"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"username": b.Username})
}

func (rt *Router) putUserPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----- messages helpers -----
type forwardBody struct {
	ConversationID string `json:"conversationId"`
}
type reactionBody struct {
	Emoji string `json:"emoji"`
}

// POST /messages/:messageId/forward
func (rt *Router) postMessageForward(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	var b forwardBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.ConversationID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"detail": "missing conversationId"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// locate source message
	msgID := ps.ByName("messageId")
	var src *Message
	for _, arr := range messages {
		for _, m := range arr {
			if m.MessageID == msgID {
				src = m
				break
			}
		}
	}
	if src == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"detail": "message not found"})
		return
	}
	// forward
	c := ensureConversation(b.ConversationID, src.Sender)
	m := &Message{
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		ConversationID: b.ConversationID,
		Sender:         src.Sender,
		Content:        src.Content,
		Type:           src.Type,
		Status:         "delivered",
		Timestamp:      time.Now().UTC(),
	}
	messages[b.ConversationID] = append(messages[b.ConversationID], m)
	c.LastMessage = m.Content
	c.Timestamp = m.Timestamp
	writeJSON(w, http.StatusCreated, m)
}

// POST /messages/:messageId/reactions
func (rt *Router) postMessageReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	var b reactionBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.Emoji == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"detail": "invalid emoji"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{
		"messageId":  ps.ByName("messageId"),
		"reactionId": uuid.Must(uuid.NewV4()).String(),
		"emoji":      b.Emoji,
	})
}

// DELETE /messages/:messageId/reactions/:reactionId
func (rt *Router) deleteMessageReaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /messages/:messageId
func (rt *Router) deleteMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----- groups (conversation-scoped) -----
type memberBody struct {
	Member string `json:"member"`
}
type nameBody struct {
	Name string `json:"name"`
}

func (rt *Router) postGroupMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	var b memberBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.Member == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"detail": "invalid member"})
		return
	}
	id := ps.ByName("conversationId")
	mu.Lock()
	c := ensureConversation(id, b.Member)
	c.Participants = append(c.Participants, b.Member)
	c.Timestamp = time.Now().UTC()
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (rt *Router) postGroupLeave(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *Router) putGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	var b nameBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"detail": "invalid name"})
		return
	}
	id := ps.ByName("conversationId")
	mu.Lock()
	ensureConversation(id).Name = b.Name
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (rt *Router) putGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, ok := getBearer(r); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"detail": "Unauthorized"})
		return
	}
	id := ps.ByName("conversationId")
	mu.Lock()
	ensureConversation(id).Photo = "uploaded"
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
