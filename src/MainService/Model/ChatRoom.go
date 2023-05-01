package Model

import "gorm.io/gorm"

type ChatRoom struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(255);" json:"name"`
	OwnerID       uint           `gorm:"type:int;not null;" json:"owner_id"`
	PosterID      uint           `gorm:"type:int;not null;" json:"poster_id"`
	Conversations []Conversation `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"conversations"`
}
