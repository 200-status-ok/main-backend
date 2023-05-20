package Model

import (
	"gorm.io/gorm"
	"time"
)

type PosterStatus string

const (
	Lost  PosterStatus = "lost"
	Found PosterStatus = "found"
)

type Poster struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Title       string         `gorm:"type:varchar(255);not null;" json:"title"`
	Description string         `gorm:"type:text" json:"description;"`
	Status      PosterStatus   `gorm:"type:status;default:'lost';not null;" json:"status"`
	UserPhone   string         `gorm:"type:varchar(15);" json:"user_phone"`
	TelegramID  string         `gorm:"type:varchar(50);" json:"telegram_id"`
	HasAlert    bool           `gorm:"type:bool;not null;default:false" json:"has_alert"`
	HasChat     bool           `gorm:"type:bool;not null;default:false" json:"has_chat"`
	Award       float64        `gorm:"type:decimal" json:"award"`
	UserID      uint           `gorm:"type:int;" json:"user_id"`
	Tags        []Tag          `gorm:"many2many:poster_tags;" json:"tags"`
	Images      []Image        `gorm:"foreignKey:PosterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"image"`
	Addresses   []Address      `gorm:"foreignKey:PosterId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"address"`
	State       string         `gorm:"type:string;default:'pending';not null;" json:"state"`
	SpecialType string         `gorm:"type:string;default:'normal';not null;" json:"special_type"`
}

func (p *Poster) GetTitle() string {
	return p.Title
}

func (p *Poster) SetTitle(title string) {
	p.Title = title
}

func (p *Poster) GetDescription() string {
	return p.Description
}

func (p *Poster) SetDescription(description string) {
	p.Description = description
}

func (p *Poster) GetStatus() PosterStatus {
	return p.Status
}

func (p *Poster) SetStatus(type_ string) {
	p.Status = PosterStatus(type_)
}

func (p *Poster) GetUserPhone() string {
	return p.UserPhone
}

func (p *Poster) SetUserPhone(userPhone string) {
	p.UserPhone = userPhone
}

func (p *Poster) GetAward() float64 {
	return p.Award
}

func (p *Poster) SetAward(award float64) {
	p.Award = award
}

func (p *Poster) GetUserID() uint {
	return p.UserID
}

func (p *Poster) SetUserID(userID uint) {
	p.UserID = userID
}

func (p *Poster) GetCategories() []Tag {
	return p.Tags
}

func (p *Poster) SetCategories(categories []Tag) {
	p.Tags = categories
}

func (p *Poster) GetImages() []Image {
	return p.Images
}

func (p *Poster) SetImages(images []Image) {
	p.Images = images
}

func (p *Poster) GetAddress() []Address {
	return p.Addresses
}

func (p *Poster) SetAddress(address []Address) {
	p.Addresses = address
}

func (p *Poster) GetTelegramID() string {
	return p.TelegramID
}

func (p *Poster) SetTelegramID(telegramID string) {
	p.TelegramID = telegramID
}

func (p *Poster) GetHasAlert() bool {
	return p.HasAlert
}

func (p *Poster) SetHasAlert(hasAlert bool) {
	p.HasAlert = hasAlert
}

func (p *Poster) GetID() uint {
	return p.ID
}

func (p *Poster) SetID(id uint) {
	p.ID = id
}

func (p *Poster) GetCreatedAt() string {
	return p.CreatedAt.String()
}

func (p *Poster) GetUpdatedAt() string {
	return p.UpdatedAt.String()
}

func (p *Poster) GetState() string {
	return p.State
}

func (p *Poster) SetState(state string) {
	p.State = state
}

func (p *Poster) GetHasChat() bool {
	return p.HasChat
}

func (p *Poster) SetHasChat(hasChat bool) {
	p.HasChat = hasChat
}

func (p *Poster) GetSpecialType() string {
	return p.SpecialType
}

func (p *Poster) SetSpecialType(specialAds string) {
	p.SpecialType = specialAds
}
