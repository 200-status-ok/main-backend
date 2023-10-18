package Model

import (
	"gorm.io/gorm"
	"time"
)

type Image struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Url       string         `gorm:"type:varchar(255)" json:"url"`
	PosterID  uint           `gorm:"not null" json:"poster_id"`
}

func (i *Image) GetUrl() string {
	return i.Url
}

func (i *Image) SetUrl(url string) {
	i.Url = url
}

func (i *Image) GetID() uint {
	return i.ID
}

func (i *Image) GetPosterID() uint {
	return i.PosterID
}

func (i *Image) SetPosterID(posterID uint) {
	i.PosterID = posterID
}
