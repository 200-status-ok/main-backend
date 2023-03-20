package Model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name    string   `gorm:"type:varchar(255);not null;unique" json:"name"`
	Posters []Poster `gorm:"many2many:poster_categories" json:"posters"`
}

func (c *Category) GetName() string {
	return c.Name
}

func (c *Category) SetName(name string) {
	c.Name = name
}

func (c *Category) GetID() uint {
	return c.ID
}

func (c *Category) GetCreatedAt() string {
	return c.CreatedAt.String()
}
