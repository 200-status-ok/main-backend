package dtos

type Message struct {
	ID             int    `json:"id"`
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	SenderID       int    `json:"sender_id"`
	ReceiverId     int    `json:"receiver_id"`
	Time           string `json:"time"`
	Type           string `json:"type"`
}

type TransferMessage struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	Type           string `json:"type"`
}
