package Model

import "gorm.io/gorm"

type PosterReport struct {
	gorm.Model
	PosterID    uint   `gorm:"not null" json:"poster_id"`
	Poster      Poster `gorm:"foreignKey:PosterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"poster"`
	IssuerID    uint   `gorm:"not null" json:"issuer_id"`
	Issuer      User   `gorm:"foreignKey:IssuerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"issuer"`
	ReportType  string `gorm:"type:varchar(50);not null;" json:"report_type"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Status      string `gorm:"type:varchar(30);not null;default:'open';" json:"status"`
}
