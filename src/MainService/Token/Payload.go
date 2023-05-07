package Token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
	UserID   uint64    `json:"userID"`
	IssuedAt time.Time `json:"issued_at"`
	Expired  time.Time `json:"expired"`
}

func NewPayload(username string, userID uint64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		Username: username,
		ID:       tokenID,
		UserID:   userID,
		IssuedAt: time.Now(),
		Expired:  time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.Expired) {
		return errors.New("token expired")
	}
	return nil
}
