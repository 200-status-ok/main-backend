package Model

import "gorm.io/gorm"

type PosterStatus string

const (
	Lost  PosterStatus = "lost"
	Found PosterStatus = "found"
)

type Poster struct {
	gorm.Model
	Title       string       `gorm:"type:varchar(255);not null;" json:"title"`
	Description string       `gorm:"type:text" json:"description;"`
	Status      PosterStatus `gorm:"type:status;default:'lost';not null;" json:"status"`
	UserPhone   string       `gorm:"type:varchar(15);" json:"user_phone"`
	TelegramID  string       `gorm:"type:varchar(50);" json:"telegram_id"`
	HasAlert    bool         `gorm:"type:bool;not null;default:false" json:"has_alert"`
	Award       float64      `gorm:"type:decimal" json:"award"`
	UserID      uint         `gorm:"type:int;" json:"user_id"`
	Categories  []Category   `gorm:"many2many:poster_categories;" json:"categories"`
	Images      []Image      `gorm:"foreignKey:PosterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"image"`
	Addresses   []Address    `gorm:"foreignKey:PosterId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"address"`
	User        User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user"`
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

func (p *Poster) GetCategories() []Category {
	return p.Categories
}

func (p *Poster) SetCategories(categories []Category) {
	p.Categories = categories
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

func (p *Poster) GetCreatedAt() string {
	return p.CreatedAt.String()
}

func (p *Poster) GetUpdatedAt() string {
	return p.UpdatedAt.String()
}
