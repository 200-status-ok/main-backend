package Model

import (
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Username  string         `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password  string         `gorm:"type:varchar(70);not null" json:"password"`
	FName     string         `gorm:"type:varchar(50);not null" json:"f_name"`
	LName     string         `gorm:"type:varchar(50);not null" json:"l_name"`
	Email     string         `gorm:"type:varchar(50);not null;unique" json:"email"`
	Phone     string         `gorm:"type:varchar(50);not null;unique" json:"phone"`
}
