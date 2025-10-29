package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Contact struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Conversation struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Members   []int64 `json:"members"`
	LastText  string  `json:"last_text"`
	UpdatedAt int64   `json:"updated_at"`
}

type ConversationDetail struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Members   []int64     `json:"members"`
	Messages  []DemoMsg   `json:"messages"`
	UpdatedAt int64       `json:"updated_at"`
}

type DemoMsg struct {
	ID        int64  `json:"id"`
	SenderID  int64  `json:"sender_id"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}

var users = []User{
	{ID: 1, Username: "mariam"},
	{ID: 2, Username: "alex"},
	{ID: 3, Username: "nino"},
}

var contacts = []Contact{
	{ID: 1, Username: "mariam"},
	{ID: 2, Username: "alex"},
	{ID: 3, Username: "nino"},
}

var conversations = []Conversation{
	{ID: 101, Name: "General chat", Members: []int64{1, 2, 3}, LastText: "hi 👋", UpdatedAt: time.Now().Unix()},
}

var convDetail = ConversationDetail{
	ID:        101,
	Name:      "General chat",
	Members:   []int64{1, 2, 3},
	UpdatedAt: time.Now().Unix(),
	Messages: []DemoMsg{
		{ID: 9001, SenderID: 2, Text: "welcome!", CreatedAt: time.Now().Add(-time.Minute * 5).Unix()},
		{ID: 9002, SenderID: 1, Text: "hi 👋", CreatedAt: time.Now().Add(-time.Minute * 4).Unix()},
	},
}

func registerDemo(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"token": "demo"})
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users" {
			usersPrefix(w, r)
			return
		}
		writeJSON(w, users)
	})
}

func usersPrefix(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/users/")
	parts := strings.Split(p, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	uidStr, tail := parts[0], strings.Join(parts[1:], "/")
	_, _ = strconv.ParseInt(uidStr, 10, 64)

	switch {
	case tail == "contacts":
		writeJSON(w, contacts)
		return
	case tail == "conversations":
		writeJSON(w, conversations)
		return
	case strings.HasPrefix(tail, "conversations/"):
		writeJSON(w, convDetail)
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
