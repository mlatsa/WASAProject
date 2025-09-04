package api

import "time"

type Message struct {
	MessageID      string    `json:"messageId"`
	ConversationID string    `json:"conversationId"`
	Sender         string    `json:"sender"`
	Content        string    `json:"content"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
}

type Conversation struct {
	ID           string    `json:"id"`
	Participants []string  `json:"participants"`
	LastMessage  string    `json:"lastMessage"`
	Timestamp    time.Time `json:"timestamp"`
	Name         string    `json:"name,omitempty"`
	Photo        string    `json:"photo,omitempty"`
}

type LoginBody struct {
	Name string `json:"name"`
}

type SendMessageInput struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}
