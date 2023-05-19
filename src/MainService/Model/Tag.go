package Model

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"type:varchar(255);not null;unique" json:"name"`
	Posters   []Poster       `gorm:"many2many:poster_tags" json:"posters"`
}

func (c *Tag) GetName() string {
	return c.Name
}

func (c *Tag) SetName(name string) {
	c.Name = name
}

func (c *Tag) GetID() uint {
	return c.ID
}

func (c *Tag) GetCreatedAt() string {
	return c.CreatedAt.String()
}
