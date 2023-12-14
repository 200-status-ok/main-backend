package Model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID             int64          `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ConversationId uint           `gorm:"not null" json:"conversation_id"`
	SenderId       uint           `gorm:"not null" json:"sender_id"`
	ReceiverId     uint           `gorm:"not null" json:"receiver_id"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	Type           string         `gorm:"type:varchar(50);not null" json:"type"`
	IsSend         bool           `gorm:"default:false" json:"is_send"`
	Status         string         `gorm:"type:varchar(50);default:'unread'" json:"status"`
}

func (m *Message) GetID() int64 {
	return m.ID
}

func (m *Message) GetConversationId() uint {
	return m.ConversationId
}

func (m *Message) SetConversationId(conversationId uint) {
	m.ConversationId = conversationId
}

func (m *Message) GetSenderId() uint {
	return m.SenderId
}

func (m *Message) SetSenderId(senderId uint) {
	m.SenderId = senderId
}

func (m *Message) GetReceiverId() uint {
	return m.ReceiverId
}

func (m *Message) SetReceiverId(receiverId uint) {
	m.ReceiverId = receiverId
}

func (m *Message) GetContent() string {
	return m.Content
}

func (m *Message) SetContent(content string) {
	m.Content = content
}
