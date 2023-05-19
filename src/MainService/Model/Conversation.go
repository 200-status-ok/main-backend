package Model

import (
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"type:varchar(255);" json:"name"`
	ImageURL  string         `gorm:"type:varchar(255);" json:"image_url"`
	OwnerID   uint           `gorm:"type:int;not null;index:idx_name,unique;" json:"owner_id"`
	MemberID  uint           `gorm:"type:int;not null;index:idx_name,unique;" json:"member_id"`
	PosterID  uint           `gorm:"type:int;not null;" json:"poster_id"`
	Messages  []Message      `gorm:"foreignKey:ConversationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"messages"`
}

func (c *Conversation) GetID() uint {
	return c.ID
}
