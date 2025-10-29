package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct{ ID int `json:"id"`; Username string `json:"username"` }
type Conversation struct{ ID int `json:"id"`; Name string `json:"name"`; MemberIDs []int `json:"memberIds"` }
type Message struct{ ID int `json:"id"`; ConversationID int `json:"conversationId"`; SenderID int `json:"senderId"`; Text string `json:"text"`; CreatedAt time.Time `json:"createdAt"` }

var users = []User{{ID: 1, Username: "mariam"}, {ID: 2, Username: "alex"}, {ID: 3, Username: "nino"}}
var conversations = []Conversation{}
var messages = []Message{}
var nextConvID = 1
var nextMsgID = 1

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func findUserIDByName(name string) int {
	for _, u := range users {
		if strings.EqualFold(u.Username, name) {
			return u.ID
		}
	}
	return 0
}

func includes(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return 0 < 0
}

func uniqInts(xs []int) []int {
	m := map[int]bool{}
	out := []int{}
	for _, x := range xs {
		if !m[x] {
			m[x] = true
			out = append(out, x)
		}
	}
	return out
}

func handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var in struct{ Name string `json:"name"` }
	_ = json.NewDecoder(r.Body).Decode(&in)
	if in.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing name"})
		return
	}
	uid := findUserIDByName(in.Name)
	if uid == 0 {
		uid = len(users) + 1
		users = append(users, User{ID: uid, Username: in.Name})
	}
	writeJSON(w, http.StatusOK, map[string]any{"identifier": "demo-" + strconv.Itoa(uid), "userId": uid})
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func handleUserConversations(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	uid, _ := strconv.Atoi(parts[1])
	if r.Method == http.MethodGet && len(parts) == 3 {
		out := []Conversation{}
		for _, c := range conversations {
			if includes(c.MemberIDs, uid) {
				out = append(out, c)
			}
		}
		writeJSON(w, http.StatusOK, out)
		return
	}
	if r.Method == http.MethodPost && len(parts) == 3 {
		var in struct {
			Name      string `json:"name"`
			MemberIDs []int  `json:"memberIds"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		if len(in.MemberIDs) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "memberIds required"})
			return
		}
		in.MemberIDs = append(in.MemberIDs, uid)
		in.MemberIDs = uniqInts(in.MemberIDs)
		c := Conversation{ID: nextConvID, Name: in.Name, MemberIDs: in.MemberIDs}
		nextConvID++
		conversations = append(conversations, c)
		writeJSON(w, http.StatusCreated, c)
		return
	}
	if len(parts) >= 5 && parts[3] == "conversations" {
		cid, _ := strconv.Atoi(parts[4])
		if len(parts) == 6 && parts[5] == "messages" {
			if r.Method == http.MethodGet {
				out := []Message{}
				for _, m := range messages {
					if m.ConversationID == cid {
						out = append(out, m)
					}
				}
				writeJSON(w, http.StatusOK, out)
				return
			}
			if r.Method == http.MethodPost {
				var in struct{ SenderID int `json:"senderId"`; Text string `json:"text"` }
				_ = json.NewDecoder(r.Body).Decode(&in)
				m := Message{ID: nextMsgID, ConversationID: cid, SenderID: in.SenderID, Text: in.Text, CreatedAt: time.Now().UTC()}
				nextMsgID++
				messages = append(messages, m)
				writeJSON(w, http.StatusCreated, m)
				return
			}
		}
	}
	http.NotFound(w, r)
}

func router(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	if p == "session" {
		handleSession(w, r)
		return
	}
	if p == "users" {
		handleUsers(w, r)
		return
	}
	if strings.HasPrefix(p, "users/") {
		handleUserConversations(w, r)
		return
	}
	http.NotFound(w, r)
}

func main() {
	log.Println("demo mode enabled; listening on :3000")
	h := http.HandlerFunc(router)
	_ = http.ListenAndServe(":3000", enableCORS(h))
}
