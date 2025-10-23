package main

import (
	"errors"
	"sync"
	"time"
)

type Message struct {
	MessageID string   `json:"messageId"`
	Sender    string   `json:"sender"`
	Content   string   `json:"content"`
	Type      string   `json:"type"`
	Timestamp int64    `json:"timestamp"`
	Reactions []string `json:"reactions,omitempty"`
}

type Conversation struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Messages    []Message `json:"messages"`
	LastMessage string    `json:"lastMessage,omitempty"`
}

type store struct {
	mu            sync.Mutex
	conversations map[string]*Conversation
	nextMsg       int
	nextReact     int
}

func newChatStore() *store {
	return &store{
		conversations: map[string]*Conversation{
			"conversation_abc": {
				ID:    "conversation_abc",
				Title: "General",
				Messages: []Message{
					{MessageID: "m1", Sender: "System", Content: "Welcome to WASAText!", Type: "text", Timestamp: time.Now().UnixMilli()},
				},
				LastMessage: "Welcome to WASAText!",
			},
		},
		nextMsg:   2,
		nextReact: 1,
	}
}

func (s *store) listConversations() []map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []map[string]string{}
	for _, c := range s.conversations {
		out = append(out, map[string]string{
			"id":          c.ID,
			"title":       c.Title,
			"lastMessage": c.LastMessage,
		})
	}
	return out
}

func (s *store) getConversation(id string) (*Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.conversations[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *c
	return &cp, nil
}

func (s *store) addMessage(convID, sender, content, mtype string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.conversations[convID]
	if !ok {
		c = &Conversation{ID: convID, Title: "New"}
		s.conversations[convID] = c
	}
	id := "m" + itoa(s.nextMsg)
	s.nextMsg++
	msg := Message{
		MessageID: id,
		Sender:    sender,
		Content:   content,
		Type:      mtype,
		Timestamp: time.Now().UnixMilli(),
	}
	c.Messages = append(c.Messages, msg)
	c.LastMessage = content
	return msg
}

func (s *store) findMessage(messageID string) (*Message, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range s.conversations {
		for i := range c.Messages {
			if c.Messages[i].MessageID == messageID {
				m := c.Messages[i]
				return &m, true
			}
		}
	}
	return nil, false
}

func (s *store) deleteMessage(messageID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range s.conversations {
		for i := range c.Messages {
			if c.Messages[i].MessageID == messageID {
				c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
				if len(c.Messages) > 0 {
					c.LastMessage = c.Messages[len(c.Messages)-1].Content
				} else {
					c.LastMessage = ""
				}
				return true
			}
		}
	}
	return false
}

func (s *store) addReaction(messageID, emoji string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range s.conversations {
		for i := range c.Messages {
			if c.Messages[i].MessageID == messageID {
				id := "r" + itoa(s.nextReact)
				s.nextReact++
				c.Messages[i].Reactions = append(c.Messages[i].Reactions, emoji+":"+id)
				return id, true
			}
		}
	}
	return "", false
}

func (s *store) removeReaction(messageID, reactionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range s.conversations {
		for i := range c.Messages {
			if c.Messages[i].MessageID == messageID {
				rs := c.Messages[i].Reactions
				for j := range rs {
					if hasSuffix(rs[j], ":"+reactionID) {
						c.Messages[i].Reactions = append(rs[:j], rs[j+1:]...)
						return true
					}
				}
			}
		}
	}
	return false
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return sign + string(b[i:])
}

func hasSuffix(s, suf string) bool {
	if len(suf) > len(s) {
		return false
	}
	return s[len(s)-len(suf):] == suf
}
