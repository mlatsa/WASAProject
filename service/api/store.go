package api

import (
	"sync"
	"time"
)

/* In-memory store */

type Store struct {
	mu            sync.Mutex
	sessions      map[string]string       // token -> username
	usernames     map[string]string       // token -> username override
	conversations map[string]*Conversation
	messages      map[string]*Message
}

func newStore() *Store {
	return &Store{
		sessions:      map[string]string{},
		usernames:     map[string]string{},
		conversations: map[string]*Conversation{},
		messages:      map[string]*Message{},
	}
}

/* Models */

type Reaction struct {
	ReactionID string `json:"reactionId"`
	Emoji      string `json:"emoji"`
}

type Message struct {
	MessageID      string     `json:"messageId"`
	ConversationID string     `json:"conversationId"`
	Sender         string     `json:"sender"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	Timestamp      time.Time  `json:"timestamp"`
	Reactions      []Reaction `json:"reactions,omitempty"`
}

type Conversation struct {
	ID           string     `json:"id"`
	Participants []string   `json:"participants"`
	Messages     []*Message `json:"messages"`
	LastMessage  string     `json:"lastMessage"`
	Timestamp    time.Time  `json:"timestamp"`
	Name         string     `json:"name,omitempty"`
	Photo        string     `json:"photo,omitempty"`
}

type ConversationDTO struct {
	ID           string     `json:"id"`
	Participants []string   `json:"participants"`
	Messages     []*Message `json:"messages,omitempty"`
	LastMessage  string     `json:"lastMessage"`
	Timestamp    time.Time  `json:"timestamp"`
	Name         string     `json:"name,omitempty"`
	Photo        string     `json:"photo,omitempty"`
}

type ConversationSummary struct {
	ID           string    `json:"id"`
	Participants []string  `json:"participants"`
	LastMessage  string    `json:"lastMessage"`
	Timestamp    time.Time `json:"timestamp"`
	Name         string    `json:"name,omitempty"`
	Photo        string    `json:"photo,omitempty"`
}

/* helpers bound to Router */

func (rt *Router) username(token string) string {
	rt.store.mu.Lock()
	defer rt.store.mu.Unlock()
	return rt.usernameLocked(token)
}

func (rt *Router) usernameLocked(token string) string {
	if u := rt.store.usernames[token]; u != "" {
		return u
	}
	if u := rt.store.sessions[token]; u != "" {
		return u
	}
	return "User"
}

func (rt *Router) ensureConversationLocked(cid, username string) *Conversation {
	if c := rt.store.conversations[cid]; c != nil {
		return c
	}
	c := &Conversation{
		ID:           cid,
		Participants: []string{username},
		Messages:     []*Message{},
		LastMessage:  "",
		Timestamp:    time.Now().UTC(),
	}
	rt.store.conversations[cid] = c
	return c
}

func (c *Conversation) toDTO() *ConversationDTO {
	return &ConversationDTO{
		ID:           c.ID,
		Participants: c.Participants,
		Messages:     c.Messages,
		LastMessage:  c.LastMessage,
		Timestamp:    c.Timestamp,
		Name:         c.Name,
		Photo:        c.Photo,
	}
}
