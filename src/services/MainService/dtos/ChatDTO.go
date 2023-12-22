package dtos

type Message struct {
	ID             int64  `json:"id"`
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	SenderID       int    `json:"sender_id"`
	ReceiverId     int    `json:"receiver_id"`
	SequenceNo     int    `json:"sequence_no"`
	Time           int64  `json:"time"`
	Type           string `json:"type"`
	Status         string `json:"status"`
}

type TransferMessage struct {
	Content        string `json:"content"`
	ConversationID int    `json:"conversation_id"`
	Type           string `json:"type"`
}
