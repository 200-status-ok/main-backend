package Model

import "gorm.io/gorm"

type MarkedPoster struct {
	gorm.Model
	UserID   uint   `gorm:"not null" json:"user_id"`
	PosterID uint   `gorm:"not null" json:"poster_id"`
	Poster   Poster `gorm:"foreignKey:PosterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"poster"`
}

func (m *MarkedPoster) GetID() uint {
	return m.ID
}

func (m *MarkedPoster) GetPosterID() uint {
	return m.PosterID
}

func (m *MarkedPoster) SetPosterID(posterID uint) {
	m.PosterID = posterID
}

func (m *MarkedPoster) GetUserID() uint {
	return m.UserID
}

func (m *MarkedPoster) SetUserID(userID uint) {
	m.UserID = userID
}

func (m *MarkedPoster) GetPoster() Poster {
	return m.Poster
}

func (m *MarkedPoster) SetPoster(poster Poster) {
	m.Poster = poster
}
