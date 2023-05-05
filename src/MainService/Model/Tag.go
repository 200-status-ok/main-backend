package Model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name    string   `gorm:"type:varchar(255);not null;unique" json:"name"`
	Posters []Poster `gorm:"many2many:poster_tags" json:"posters"`
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
