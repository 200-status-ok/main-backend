package Model

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Amount    float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	TrackID   string         `gorm:"not null" json:"track_id"`
	Status    string         `gorm:"type:varchar(50);not null" json:"status"`
}

func (p *Payment) GetID() uint {
	return p.ID
}

func (p *Payment) SetID(id uint) {
	p.ID = id
}

func (p *Payment) GetCreatedAt() string {
	return p.CreatedAt.String()
}

func (p *Payment) GetUpdatedAt() string {
	return p.UpdatedAt.String()
}

func (p *Payment) GetUserID() uint {
	return p.UserID
}

func (p *Payment) SetUserID(userID uint) {
	p.UserID = userID
}

func (p *Payment) GetAmount() float64 {
	return p.Amount
}

func (p *Payment) SetAmount(amount float64) {
	p.Amount = amount
}

func (p *Payment) GetTrackID() string {
	return p.TrackID
}

func (p *Payment) SetTrackID(trackID string) {
	p.TrackID = trackID
}

func (p *Payment) GetStatus() string {
	return p.Status
}

func (p *Payment) SetStatus(status string) {
	p.Status = status
}

func (p *Payment) GetPayment() *Payment {
	return p
}

func (p *Payment) SetPayment(payment *Payment) {
	p.ID = payment.ID
	p.UserID = payment.UserID
	p.Amount = payment.Amount
	p.TrackID = payment.TrackID
	p.Status = payment.Status
}
