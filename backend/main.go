package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

/* ===================== Types & Storage ===================== */

type Message struct {
	MessageID      string  `json:"messageId"`
	ConversationID string  `json:"conversationId"`
	Sender         string  `json:"sender"`
	Content        string  `json:"content"`
	Type           string  `json:"type"`   // e.g. "text"
	Status         string  `json:"status"` // e.g. "delivered"
	Timestamp      float64 `json:"timestamp"`
}

type Conversation struct {
	ID           string   `json:"id"`
	Participants []string `json:"participants"`
	LastMessage  string   `json:"lastMessage,omitempty"`
	Timestamp    float64  `json:"timestamp"`
	// optional fields
	Name  string `json:"name,omitempty"`
	Photo string `json:"photo,omitempty"`
}

var (
	// auth
	usersByName   = map[string]string{} // name -> identifier
	usernamesByID = map[string]string{} // identifier -> name

	// data
	conversations     = map[string]*Conversation{}     // id -> conversation
	messagesByConvo   = map[string][]*Message{}        // convo -> messages
	reactionsByMsg    = map[string]map[string]string{} // msgID -> (reactionID -> value)
	userPhotos        = map[string]string{}            // userID -> url

	mu sync.Mutex
)

/* ===================== Helpers ===================== */

func now() float64 { return float64(time.Now().UnixMilli()) / 1000.0 }

func randID(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func readJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

/* ===================== CORS ===================== */

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow-all as per brief; Max-Age=1 on preflight
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		w.Header().Set("Access-Control-Max-Age", "1")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/* ===================== Auth ===================== */

func bearerID(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	id := strings.TrimSpace(parts[1])
	mu.Lock()
	defer mu.Unlock()
	_, ok := usernamesByID[id]
	return id, ok
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if id, ok := bearerID(r); ok {
			// inject into context via header (simple)
			r.Header.Set("X-User-ID", id)
			next(w, r)
			return
		}
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
}

/* ===================== Route helpers ===================== */

func pathAfter(r *http.Request, prefix string) (string, bool) {
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return "", false
	}
	return strings.TrimPrefix(r.URL.Path, prefix), true
}

/* ===================== Handlers ===================== */

// GET /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /session  (doLogin)
func sessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}
	if l := len(body.Name); l < 3 || l > 16 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid username"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	id, exists := usersByName[body.Name]
	if !exists {
		id = randID(24)
		usersByName[body.Name] = id
	}
	// ensure reverse map is set (survives re-login after restart)
	usernamesByID[id] = body.Name

	writeJSON(w, http.StatusCreated, map[string]string{"identifier": id})
}

// POST /user/username  (setMyUserName)
func setMyUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}
	if l := len(body.Name); l < 3 || l > 16 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid username"})
		return
	}

	userID := r.Header.Get("X-User-ID")
	mu.Lock()
	defer mu.Unlock()

	// name taken by someone else?
	if existing, ok := usersByName[body.Name]; ok && existing != userID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid username"})
		return
	}
	// remove old mapping (if any)
	if oldName, ok := usernamesByID[userID]; ok {
		if mapped, ok2 := usersByName[oldName]; ok2 && mapped == userID {
			delete(usersByName, oldName)
		}
	}
	// set new
	usersByName[body.Name] = userID
	usernamesByID[userID] = body.Name

	writeJSON(w, http.StatusOK, map[string]string{"message": "OK"})
}

// GET /conversations  (getMyConversations)
func getMyConversationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	userID := r.Header.Get("X-User-ID")

	mu.Lock()
	defer mu.Unlock()

	var mine []Conversation
	for _, c := range conversations {
		for _, p := range c.Participants {
			if p == userID {
				mine = append(mine, *c)
				break
			}
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"conversations": mine})
}

// GET /conversations/{conversationId}  (getConversation)
func getConversationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	userID := r.Header.Get("X-User-ID")
	tail, _ := pathAfter(r, "/conversations/")
	convID := strings.Trim(tail, "/")

	mu.Lock()
	defer mu.Unlock()

	conv, ok := conversations[convID]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Conversation not found"})
		return
	}
	// must participate
	part := false
	for _, p := range conv.Participants {
		if p == userID {
			part = true
			break
		}
	}
	if !part {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Conversation not found"})
		return
	}

	msgs := messagesByConvo[convID]
	writeJSON(w, http.StatusOK, map[string]any{
		"conversation": map[string]any{
			"id":           conv.ID,
			"participants": conv.Participants,
			"lastMessage":  conv.LastMessage,
			"timestamp":    conv.Timestamp,
			"name":         conv.Name,
			"messages":     msgs,
		},
	})
}

// POST /conversations/{conversationId}/messages  (sendMessage)
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	userID := r.Header.Get("X-User-ID")

	// extract convo id
	tail, _ := pathAfter(r, "/conversations/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	if len(parts) < 2 || parts[1] != "messages" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	convID := parts[0]

	var body struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}
	if err := readJSON(r, &body); err != nil || strings.TrimSpace(body.Content) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}
	if body.Type == "" {
		body.Type = "text"
	}

	mu.Lock()
	defer mu.Unlock()

	// ensure conversation exists (create if missing with sender as participant)
	conv, ok := conversations[convID]
	if !ok {
		conv = &Conversation{
			ID:           convID,
			Participants: []string{userID},
			Timestamp:    now(),
		}
		conversations[convID] = conv
	}
	// ensure sender is participant
	isPart := false
	for _, p := range conv.Participants {
		if p == userID {
			isPart = true
			break
		}
	}
	if !isPart {
		conv.Participants = append(conv.Participants, userID)
	}

	msg := &Message{
		MessageID:      "msg_" + randID(12),
		ConversationID: convID,
		Sender:         userID,
		Content:        body.Content,
		Type:           body.Type,
		Status:         "delivered",
		Timestamp:      now(),
	}
	messagesByConvo[convID] = append(messagesByConvo[convID], msg)
	conv.LastMessage = body.Content
	conv.Timestamp = now()

	writeJSON(w, http.StatusCreated, msg)
}

// POST /messages/{messageId}/forward  (forwardMessage)
func forwardMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	userID := r.Header.Get("X-User-ID")

	tail, _ := pathAfter(r, "/messages/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	if len(parts) < 2 || parts[1] != "forward" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	msgID := parts[0]

	var body struct {
		ConversationID string `json:"conversationId"`
	}
	if err := readJSON(r, &body); err != nil || strings.TrimSpace(body.ConversationID) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// find source message
	var src *Message
	for _, msgs := range messagesByConvo {
		for _, m := range msgs {
			if m.MessageID == msgID {
				src = m
				break
			}
		}
		if src != nil {
			break
		}
	}
	if src == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Message not found"})
		return
	}

	// ensure destination conversation
	conv, ok := conversations[body.ConversationID]
	if !ok {
		conv = &Conversation{
			ID:           body.ConversationID,
			Participants: []string{userID},
			Timestamp:    now(),
		}
		conversations[body.ConversationID] = conv
	}
	// ensure caller is participant
	part := false
	for _, p := range conv.Participants {
		if p == userID {
			part = true
			break
		}
	}
	if !part {
		conv.Participants = append(conv.Participants, userID)
	}

	fwd := &Message{
		MessageID:      "msg_" + randID(12),
		ConversationID: body.ConversationID,
		Sender:         userID,
		Content:        src.Content,
		Type:           src.Type,
		Status:         "delivered",
		Timestamp:      now(),
	}
	messagesByConvo[body.ConversationID] = append(messagesByConvo[body.ConversationID], fwd)
	conv.LastMessage = fwd.Content
	conv.Timestamp = now()

	writeJSON(w, http.StatusOK, fwd)
}

// POST /messages/{messageId}/reactions  (commentMessage)
func addReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/messages/")
	msgID := strings.TrimSuffix(strings.Trim(tail, "/"), "/reactions")

	var body struct {
		Reaction string `json:"reaction"`
	}
	if err := readJSON(r, &body); err != nil || strings.TrimSpace(body.Reaction) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// ensure message exists
	found := false
	for _, msgs := range messagesByConvo {
		for _, m := range msgs {
			if m.MessageID == msgID {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Message not found"})
		return
	}

	rid := "react_" + randID(10)
	if _, ok := reactionsByMsg[msgID]; !ok {
		reactionsByMsg[msgID] = map[string]string{}
	}
	reactionsByMsg[msgID][rid] = body.Reaction

	writeJSON(w, http.StatusOK, map[string]string{"reactionId": rid})
}

// DELETE /messages/{messageId}/reactions/{reactionId}  (uncommentMessage)
func removeReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/messages/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	// {messageId}/reactions/{reactionId}
	if len(parts) != 3 || parts[1] != "reactions" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	msgID, reactionID := parts[0], parts[2]

	mu.Lock()
	defer mu.Unlock()

	if store, ok := reactionsByMsg[msgID]; ok {
		if _, ok2 := store[reactionID]; ok2 {
			delete(store, reactionID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Reaction not found"})
}

// DELETE /messages/{messageId}  (deleteMessage)
func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/messages/")
	msgID := strings.Trim(tail, "/")

	mu.Lock()
	defer mu.Unlock()

	for convID, msgs := range messagesByConvo {
		for i, m := range msgs {
			if m.MessageID == msgID {
				messagesByConvo[convID] = append(msgs[:i], msgs[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Message not found"})
}

// POST /groups/{conversationId}/members  (addToGroup)
func addToGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/groups/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	// {conversationId}/members
	if len(parts) != 2 || parts[1] != "members" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	convID := parts[0]

	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil || strings.TrimSpace(body.ID) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	conv, ok := conversations[convID]
	if !ok {
		conv = &Conversation{
			ID:           convID,
			Participants: []string{},
			Timestamp:    now(),
		}
		conversations[convID] = conv
	}
	// add if missing
	found := false
	for _, p := range conv.Participants {
		if p == body.ID {
			found = true
			break
		}
	}
	if !found {
		conv.Participants = append(conv.Participants, body.ID)
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "OK"})
}

// POST /groups/{conversationId}/leave  (leaveGroup)
func leaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/groups/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	// {conversationId}/leave
	if len(parts) != 2 || parts[1] != "leave" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	convID := parts[0]
	userID := r.Header.Get("X-User-ID")

	mu.Lock()
	defer mu.Unlock()

	conv, ok := conversations[convID]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Conversation not found"})
		return
	}
	// remove user if present
	out := make([]string, 0, len(conv.Participants))
	for _, p := range conv.Participants {
		if p != userID {
			out = append(out, p)
		}
	}
	conv.Participants = out
	writeJSON(w, http.StatusOK, map[string]string{"message": "OK"})
}

// POST /groups/{conversationId}/name  (setGroupName)
func setGroupNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/groups/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	// {conversationId}/name
	if len(parts) != 2 || parts[1] != "name" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	convID := parts[0]

	var body struct {
		Name string `json:"name"`
	}
	if err := readJSON(r, &body); err != nil || len(strings.TrimSpace(body.Name)) < 3 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid name"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	conv, ok := conversations[convID]
	if !ok {
		conv = &Conversation{
			ID:           convID,
			Participants: []string{},
			Timestamp:    now(),
		}
		conversations[convID] = conv
	}
	conv.Name = body.Name
	writeJSON(w, http.StatusOK, map[string]string{"message": "OK"})
}

// POST /user/photo  (setMyPhoto)
func setMyPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	var body struct {
		MediaURL string `json:"mediaUrl"`
	}
	if err := readJSON(r, &body); err != nil || strings.TrimSpace(body.MediaURL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}
	userID := r.Header.Get("X-User-ID")

	mu.Lock()
	userPhotos[userID] = body.MediaURL
	mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"message": "OK"})
}

// POST /groups/{conversationId}/photo  (setGroupPhoto)
func setGroupPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	tail, _ := pathAfter(r, "/groups/")
	parts := strings.Split(strings.Trim(tail, "/"), "/")
	// {conversationId}/photo
	if len(parts) != 2 || parts[1] != "photo" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	convID := parts[0]

	var body struct {
		MediaURL string `json:"mediaUrl"`
	}
	if err := readJSON(r, &body); err != nil || strings.TrimSpace(body.MediaURL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	conv, ok := conversations[convID]
	if !ok {
		conv = &Conversation{
			ID:           convID,
			Participants: []string{},
			Timestamp:    now(),
		}
		conversations[convID] = conv
	}
	conv.Photo = body.MediaURL
	writeJSON(w, http.StatusOK, map[string]string{"message": "OK"})
}

/* ===================== main ===================== */

func main() {
	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()

	// public
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/session", sessionHandler)

	// protected
	mux.Handle("/user/username", withCORS(requireAuth(setMyUsernameHandler)))
	mux.Handle("/conversations", withCORS(requireAuth(getMyConversationsHandler)))
	mux.Handle("/conversations/", withCORS(requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/messages") && r.Method == http.MethodPost {
			sendMessageHandler(w, r)
			return
		}
		getConversationHandler(w, r)
	})))
	mux.Handle("/messages/", withCORS(requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/forward") && r.Method == http.MethodPost {
			forwardMessageHandler(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "/reactions/") && r.Method == http.MethodDelete {
			removeReactionHandler(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/reactions") && r.Method == http.MethodPost {
			addReactionHandler(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			deleteMessageHandler(w, r)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
	})))
	mux.Handle("/groups/", withCORS(requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/members") && r.Method == http.MethodPost {
			addToGroupHandler(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/leave") && r.Method == http.MethodPost {
			leaveGroupHandler(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/name") && r.Method == http.MethodPost {
			setGroupNameHandler(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/photo") && r.Method == http.MethodPost {
			setGroupPhotoHandler(w, r)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
	})))
	mux.Handle("/user/photo", withCORS(requireAuth(setMyPhotoHandler)))

	// wrap top-level with CORS as well (covers public routes & OPTIONS)
	handler := withCORS(mux)

	log.Println("Go server on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
