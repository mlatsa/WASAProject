package api

import (
	"sync"
	"time"
)

var (
	mu            sync.Mutex
	users         = map[string]string{}        // identifier -> name
	conversations = map[string]*Conversation{} // id -> conversation
	messages      = map[string][]*Message{}    // conversationID -> messages
)

func ensureConversation(id string, participants ...string) *Conversation {
	c, ok := conversations[id]
	if !ok {
		c = &Conversation{
			ID:           id,
			Participants: participants,
			LastMessage:  "",
			Timestamp:    time.Now().UTC(),
		}
		conversations[id] = c
	}
	return c
}
