package Repository

import (
	"errors"
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*Model.User, error) {
	var user Model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) UserUpdate(user *Model.User) error {
	err := r.db.Model(&user).Updates(user).Error
	if err != nil {
		return errors.New("error while updating user")
	}
	return nil
}

func (r *UserRepository) UserCreate(user *Model.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return errors.New("error while creating user")
	}
	return nil
}
