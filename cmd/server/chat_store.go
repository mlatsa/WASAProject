package main

import (
	"errors"
	"sync"
	"time"
)

type Message struct {
	MessageID string            `json:"messageId"`
	Sender    string            `json:"sender"`
	Content   string            `json:"content"`
	Type      string            `json:"type"`
	Timestamp int64             `json:"timestamp"`
	Reactions map[string]string `json:"reactions,omitempty"` // reactionId -> emoji
}

type Conversation struct {
	ID          string     `json:"id"`
	Title       string     `json:"title,omitempty"`
	Messages    []*Message `json:"messages,omitempty"`
	LastMessage string     `json:"lastMessage,omitempty"`
}

type ChatStore struct {
	mu        sync.Mutex
	convs     map[string]*Conversation
	nextMsgID int
	nextReact int
}

func newChatStore() *ChatStore {
	cs := &ChatStore{
		convs:     make(map[string]*Conversation),
		nextMsgID: 1,
		nextReact: 1,
	}
	cs.convs["conversation_abc"] = &Conversation{
		ID:    "conversation_abc",
		Title: "General",
	}
	cs.addMessage("conversation_abc", "System", "Welcome to WASAText!", "text")
	return cs
}

func (cs *ChatStore) listConversations() []*Conversation {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	out := make([]*Conversation, 0, len(cs.convs))
	for _, c := range cs.convs {
		if n := len(c.Messages); n > 0 {
			c.LastMessage = c.Messages[n-1].Content
		} else {
			c.LastMessage = ""
		}
		out = append(out, &Conversation{ID: c.ID, Title: c.Title, LastMessage: c.LastMessage})
	}
	return out
}

func (cs *ChatStore) getConversation(id string) (*Conversation, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c, ok := cs.convs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copyC := &Conversation{ID: c.ID, Title: c.Title, LastMessage: c.LastMessage}
	copyC.Messages = append([]*Message(nil), c.Messages...)
	return copyC, nil
}

func (cs *ChatStore) ensureConv(id string) *Conversation {
	if c, ok := cs.convs[id]; ok {
		return c
	}
	c := &Conversation{ID: id, Title: id}
	cs.convs[id] = c
	return c
}

func (cs *ChatStore) addMessage(convID, sender, content, typ string) *Message {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c := cs.ensureConv(convID)
	id := cs.nextMsgID
	cs.nextMsgID++
	m := &Message{
		MessageID: "m" + itoa(id),
		Sender:    sender,
		Content:   content,
		Type:      typ,
		Timestamp: time.Now().UnixMilli(),
		Reactions: map[string]string{},
	}
	c.Messages = append(c.Messages, m)
	c.LastMessage = content
	return m
}

func (cs *ChatStore) deleteMessage(msgID string) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.convs {
		for i, m := range c.Messages {
			if m.MessageID == msgID {
				c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
				if n := len(c.Messages); n > 0 {
					c.LastMessage = c.Messages[n-1].Content
				} else {
					c.LastMessage = ""
				}
				return true
			}
		}
	}
	return false
}

func (cs *ChatStore) findMessage(msgID string) (*Message, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.convs {
		for _, m := range c.Messages {
			if m.MessageID == msgID {
				return m, true
			}
		}
	}
	return nil, false
}

func (cs *ChatStore) addReaction(msgID, emoji string) (string, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.convs {
		for _, m := range c.Messages {
			if m.MessageID == msgID {
				rid := "r" + itoa(cs.nextReact)
				cs.nextReact++
				if m.Reactions == nil {
					m.Reactions = map[string]string{}
				}
				m.Reactions[rid] = emoji
				return rid, true
			}
		}
	}
	return "", false
}

func (cs *ChatStore) removeReaction(msgID, reactionID string) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.convs {
		for _, m := range c.Messages {
			if m.MessageID == msgID {
				if m.Reactions != nil {
					if _, ok := m.Reactions[reactionID]; ok {
						delete(m.Reactions, reactionID)
						return true
					}
				}
			}
		}
	}
	return false
}

// minimal int->string helper to avoid extra imports
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + (i % 10))
		i /= 10
	}
	return string(buf[pos:])
}
