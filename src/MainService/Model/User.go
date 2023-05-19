package Model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                  uint           `gorm:"primarykey" json:"id"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Username            string         `gorm:"type:varchar(50);not null;unique" json:"username"`
	Posters             []Poster       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"posters"`
	Wallet              float64        `gorm:"type:decimal(10,2);default:0" json:"wallet"`
	OwnConversations    []Conversation `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"own_conversations"`
	MemberConversations []Conversation `gorm:"foreignKey:MemberID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"member_conversations"`
	MarkedPosters       []MarkedPoster `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"marked_posters"`
	Payments            []Payment      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"Payments"`
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

func (u *User) SetWallet(wallet float64) {
	u.Wallet = wallet
}

func (u *User) GetWallet() float64 {
	return u.Wallet
}
