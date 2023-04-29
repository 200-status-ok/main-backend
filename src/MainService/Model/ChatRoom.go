package Model

import "gorm.io/gorm"

type ChatRoom struct {
	gorm.Model
	Owner         uint           `gorm:"type:int;not null;" json:"owner"`
	PosterID      uint           `gorm:"type:int;not null;" json:"poster_id"`
	Conversations []Conversation `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"conversations"`
}
