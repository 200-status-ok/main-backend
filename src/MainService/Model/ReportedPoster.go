package Model

import (
	"gorm.io/gorm"
	"time"
)

type PosterReport struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	PosterID    uint           `gorm:"not null" json:"poster_id"`
	Poster      Poster         `gorm:"foreignKey:PosterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"poster"`
	IssuerID    uint           `gorm:"not null" json:"issuer_id"`
	Issuer      User           `gorm:"foreignKey:IssuerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"issuer"`
	ReportType  string         `gorm:"type:varchar(50);not null;" json:"report_type"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	Status      string         `gorm:"type:varchar(30);not null;default:'open';" json:"status"`
}
