package dtos

import "time"

type Message struct {
	ID             int64     `json:"id"`
	Content        string    `json:"content"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	ReceiverId     int       `json:"receiver_id"`
	Time           time.Time `json:"time"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
}

type TransferMessage struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	Type           string `json:"type"`
}
