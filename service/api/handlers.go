package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func getBearer(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	const p = "Bearer "
	if len(h) > len(p) && h[:len(p)] == p {
		return h[len(p):], true
	}
	return "", false
}

func (rt *Router) getHelloWorld(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "hello"})
}

func (rt *Router) liveness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /session  { "name": "Alex" }
func (rt *Router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body LoginBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Name) < 3 {
		http.Error(w, `{"error":"invalid name"}`, http.StatusBadRequest)
		return
	}
	id := uuid.Must(uuid.NewV4()).String()
	mu.Lock()
	users[id] = body.Name
	mu.Unlock()
	writeJSON(w, http.StatusCreated, map[string]string{"identifier": id})
}

// GET /conversations
func (rt *Router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id, ok := getBearer(r)
	if !ok || users[id] == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	list := make([]*Conversation, 0, len(conversations))
	for _, c := range conversations {
		list = append(list, c)
	}
	writeJSON(w, http.StatusOK, map[string]any{"conversations": list})
}

// GET /conversations/:id
func (rt *Router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, ok := getBearer(r)
	if !ok || users[user] == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	cid := ps.ByName("id")
	mu.Lock()
	defer mu.Unlock()
	c := ensureConversation(cid, users[user])
	writeJSON(w, http.StatusOK, map[string]any{
		"conversation": map[string]any{
			"id":           c.ID,
			"participants": c.Participants,
			"lastMessage":  c.LastMessage,
			"timestamp":    c.Timestamp,
			"messages":     messages[cid],
		},
	})
}

// POST /conversations/:id/messages  { "content": "hey!", "type": "text" }
func (rt *Router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, ok := getBearer(r)
	if !ok || users[user] == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	cid := ps.ByName("id")

	var body SendMessageInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Content == "" {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	c := ensureConversation(cid, users[user])

	m := &Message{
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		ConversationID: cid,
		Sender:         users[user],
		Content:        body.Content,
		Type:           ifEmpty(body.Type, "text"),
		Status:         "delivered",
		Timestamp:      time.Now().UTC(),
	}
	messages[cid] = append(messages[cid], m)
	c.LastMessage = m.Content
	c.Timestamp = m.Timestamp

	writeJSON(w, http.StatusCreated, m)
}

func ifEmpty(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
