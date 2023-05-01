package Model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username      string         `gorm:"type:varchar(50);not null;unique" json:"username"`
	Posters       []Poster       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"posters"`
	ChatRooms     []ChatRoom     `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"chat_rooms"`
	Conversations []Conversation `gorm:"foreignKey:MemberID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"conversations"`
	MarkedPosters []MarkedPoster `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"marked_posters"`
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) SetID(id uint) {
	u.ID = id
}

func (u *User) GetCreatedAt() string {
	return u.CreatedAt.String()
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.CreatedAt = createdAt
}

func (u *User) GetUpdatedAt() string {
	return u.UpdatedAt.String()
}

func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.UpdatedAt = updatedAt
}
