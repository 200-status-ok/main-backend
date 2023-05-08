package Model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password string `gorm:"type:varchar(70);not null" json:"password"`
	FName    string `gorm:"type:varchar(50);not null" json:"f_name"`
	LName    string `gorm:"type:varchar(50);not null" json:"l_name"`
	Email    string `gorm:"type:varchar(50);not null;unique" json:"email"`
	Phone    string `gorm:"type:varchar(50);not null;unique" json:"phone"`
}
