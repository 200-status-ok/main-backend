package Model

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Url      string `gorm:"type:varchar(255);not null" json:"url"`
	PosterID uint   `gorm:"not null" json:"image_id"`
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
