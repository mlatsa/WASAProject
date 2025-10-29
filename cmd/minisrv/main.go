package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
type Message struct {
	ID             int64  `json:"id"`
	ConversationID int64  `json:"conversation_id"`
	SenderID       int64  `json:"sender_id"`
	Text           string `json:"text"`
	CreatedAt      int64  `json:"created_at"`
}
type Conversation struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Members   []int64   `json:"members"`
	LastText  string    `json:"last_text"`
	UpdatedAt int64     `json:"updated_at"`
	Messages  []Message `json:"messages,omitempty"`
}

var users = []User{
	{ID: 1, Username: "mariam"},
	{ID: 2, Username: "alex"},
	{ID: 3, Username: "nino"},
}
var contacts = map[int64][]int64{
	1: {2, 3},
}
var convs = []Conversation{
	{ID: 101, Name: "General chat", Members: []int64{1, 2, 3}, LastText: "hi 👋", UpdatedAt: time.Now().Unix()},
}
var nextConvID int64 = 200

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
func parseInt(s string) int64 { n, _ := strconv.ParseInt(s, 10, 64); return n }

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})

	// POST /session {name}
	// Return both "userId" and "user_id" + a JWT-shaped token to satisfy jwt-decode.
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var b struct{ Name string `json:"name"` }
		_ = json.NewDecoder(r.Body).Decode(&b)
		uid := int64(1)
		username := "mariam"
		for _, u := range users {
			if strings.EqualFold(u.Username, b.Name) {
				uid = u.ID
				username = u.Username
				break
			}
		}
		// header.payload.signature (payload: {"uid":uid})
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
			"eyJ1aWQiOi" + strconv.FormatInt(uid, 10) + "fQ." +
			"fakeSig"
		writeJSON(w, map[string]interface{}{
			"token":    token,
			"userId":   uid,       // camelCase for your LoginView
			"user_id":  uid,       // snake_case just in case
			"username": username,  // handy
		})
	})

	// Fallback GET /conversations for some UIs that call it without /users/:id
	mux.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			writeJSON(w, convs)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	// GET /users (with optional ?name= search)
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		q := strings.TrimSpace(r.URL.Query().Get("name"))
		if q == "" {
			writeJSON(w, users)
			return
		}
		out := []User{}
		for _, u := range users {
			if strings.Contains(strings.ToLower(u.Username), strings.ToLower(q)) {
				out = append(out, u)
			}
		}
		writeJSON(w, out)
	})

	// /users/:uid/...
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		parts := strings.Split(path, "/")

		// /users/:uid/contacts
		if len(parts) >= 2 && parts[1] == "contacts" {
			uid := parseInt(parts[0])
			switch r.Method {
			case http.MethodGet:
				ids := contacts[uid]
				out := []User{}
				for _, id := range ids {
					for _, u := range users {
						if u.ID == id {
							out = append(out, u)
						}
					}
				}
				writeJSON(w, out)
				return
			case http.MethodPost:
				var b struct{ ContactID int64 `json:"contactId"` }
				_ = json.NewDecoder(r.Body).Decode(&b)
				contacts[uid] = append(contacts[uid], b.ContactID)
				w.WriteHeader(http.StatusCreated)
				writeJSON(w, map[string]string{"status": "added"})
				return
			case http.MethodDelete:
				if len(parts) == 3 {
					cid := parseInt(parts[2])
					cur := contacts[uid]
					nw := make([]int64, 0, len(cur))
					for _, x := range cur {
						if x != cid {
							nw = append(nw, x)
						}
					}
					contacts[uid] = nw
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// /users/:uid/conversations ...
		if len(parts) >= 2 && parts[1] == "conversations" {
			uid := parseInt(parts[0])
			if len(parts) == 2 {
				// /users/:uid/conversations
				switch r.Method {
				case http.MethodGet:
					writeJSON(w, convs)
					return
				case http.MethodPost:
					var b struct {
						Name    string  `json:"name"`
						Members []int64 `json:"members"`
					}
					_ = json.NewDecoder(r.Body).Decode(&b)
					nextConvID++
					c := Conversation{ID: nextConvID, Name: b.Name, Members: b.Members, UpdatedAt: time.Now().Unix()}
					if len(c.Members) == 0 {
						c.Members = []int64{uid}
					}
					convs = append(convs, c)
					writeJSON(w, c)
					return
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			// /users/:uid/conversations/:cid[/...]
			if len(parts) >= 3 {
				cid := parseInt(parts[2])
				var c *Conversation
				for i := range convs {
					if convs[i].ID == cid {
						c = &convs[i]
						break
					}
				}
				if c == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				if len(parts) == 3 && r.Method == http.MethodGet {
					writeJSON(w, c)
					return
				}
				// messages
				if len(parts) >= 4 && parts[3] == "messages" {
					switch r.Method {
					case http.MethodGet:
						writeJSON(w, c.Messages)
						return
					case http.MethodPost:
						var b struct {
							Text string `json:"text"`
						}
						_ = json.NewDecoder(r.Body).Decode(&b)
						m := Message{
							ID:             time.Now().UnixNano(),
							ConversationID: cid,
							SenderID:       uid,
							Text:           b.Text,
							CreatedAt:      time.Now().Unix(),
						}
						c.Messages = append(c.Messages, m)
						c.LastText = b.Text
						c.UpdatedAt = time.Now().Unix()
						writeJSON(w, m)
						return
					case http.MethodDelete:
						// ignore deletes gracefully
						w.WriteHeader(http.StatusNoContent)
						return
					}
				}
				// reactions endpoints: accept/no-op
				reReact := regexp.MustCompile(`^messages/\d+/reaction`)
				if len(parts) >= 5 && parts[3] == "messages" && reReact.MatchString(strings.Join(parts[4:], "/")) {
					if r.Method == http.MethodPost || r.Method == http.MethodDelete {
						w.WriteHeader(http.StatusNoContent)
						return
					}
				}
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	addr := ":3000"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}
