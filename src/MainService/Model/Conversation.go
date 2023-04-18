package Model

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	User1Id  uint      `gorm:"not null" json:"user_1_id"`
	User2Id  uint      `gorm:"not null" json:"user_2_id"`
	Messages []Message `gorm:"foreignKey:ConversationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"messages"`
}

func (c *Conversation) GetID() uint {
	return c.ID
}

func (c *Conversation) GetUser1Id() uint {
	return c.User1Id
}

func (c *Conversation) SetUser1Id(user1Id uint) {
	c.User1Id = user1Id
}

func (c *Conversation) GetUser2Id() uint {
	return c.User2Id
}

func (c *Conversation) SetUser2Id(user2Id uint) {
	c.User2Id = user2Id
}
