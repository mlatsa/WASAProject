package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

var chatStore = newChatStore()

// Register real chat endpoints consumed by the UI.
func (rt *Router) registerChatAPI() {
	// GET /api/conversations
	rt.mux.HandleFunc("/api/conversations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"conversations": chatStore.listConversations(),
		})
	})

	// GET /api/conversations/{id}
	rt.mux.HandleFunc("/api/conversations/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r); return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/conversations/")
		if strings.Contains(id, "/") || id == "" {
			http.NotFound(w, r); return
		}
		conv, err := chatStore.getConversation(id)
		if err != nil {
			http.NotFound(w, r); return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"conversation": conv})
	})

	// POST /api/conversations/{id}/messages
	rt.mux.HandleFunc("/api/conversations/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/conversations/") {
			http.NotFound(w, r); return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/conversations/")
		parts := strings.Split(rest, "/")
		if len(parts) == 2 && parts[1] == "messages" && r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			var payload struct {
				Content string `json:"content"`
				Type    string `json:"type"`
			}
			_ = json.Unmarshal(body, &payload)
			if payload.Type == "" { payload.Type = "text" }
			msg := chatStore.addMessage(parts[0], "Alex", payload.Content, payload.Type)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": msg.MessageID})
			return
		}
	})

	// DELETE /api/messages/{messageId}
	rt.mux.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/messages/") {
			http.NotFound(w, r); return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/messages/")
		parts := strings.Split(rest, "/")
		if len(parts) == 1 && r.Method == http.MethodDelete {
			if ok := chatStore.deleteMessage(parts[0]); ok {
				w.WriteHeader(http.StatusNoContent)
			} else {
				http.NotFound(w, r)
			}
			return
		}
	})

	// POST /api/messages/{messageId}/reactions
	rt.mux.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/messages/") {
			http.NotFound(w, r); return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/messages/")
		parts := strings.Split(rest, "/")
		if len(parts) == 2 && parts[1] == "reactions" && r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			var payload struct{ Emoji string `json:"emoji"` }
			_ = json.Unmarshal(body, &payload)
			if payload.Emoji == "" { payload.Emoji = "👍" }
			if rid, ok := chatStore.addReaction(parts[0], payload.Emoji); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"reactionId": rid})
			} else {
				http.NotFound(w, r)
			}
			return
		}
	})

	// DELETE /api/messages/{messageId}/reactions/{reactionId}
	rt.mux.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/messages/") {
			http.NotFound(w, r); return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/messages/")
		parts := strings.Split(rest, "/")
		if len(parts) == 3 && parts[1] == "reactions" && r.Method == http.MethodDelete {
			if ok := chatStore.removeReaction(parts[0], parts[2]); ok {
				w.WriteHeader(http.StatusNoContent)
			} else {
				http.NotFound(w, r)
			}
			return
		}
	})

	// POST /api/messages/{messageId}/forward
	rt.mux.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/messages/") {
			http.NotFound(w, r); return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/messages/")
		parts := strings.Split(rest, "/")
		if len(parts) == 2 && parts[1] == "forward" && r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			var payload struct{ ConversationID string `json:"conversationId"` }
			_ = json.Unmarshal(body, &payload)
			if payload.ConversationID == "" {
				http.Error(w, "conversationId required", http.StatusBadRequest); return
			}
			if msg, ok := chatStore.findMessage(parts[0]); ok {
				newMsg := chatStore.addMessage(payload.ConversationID, msg.Sender, msg.Content, msg.Type)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": newMsg.MessageID})
				return
			}
			http.NotFound(w, r)
			return
		}
	})
}
