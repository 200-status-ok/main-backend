package dtos

type Message struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	SenderID       int    `json:"sender"`
	ReceiverId     int    `json:"receiver"`
	Time           string `json:"time"`
	Type           string `json:"type"`
}

type TransferMessage struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	Type           string `json:"type"`
}
