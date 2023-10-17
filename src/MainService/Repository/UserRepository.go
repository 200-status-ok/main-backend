package Repository

import (
	"errors"
	"github.com/200-status-ok/main-backend/src/MainService/Model"
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
	err := r.db.Unscoped().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.DeletedAt.Valid {
		r.db.Unscoped().Delete(&user)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) UserUpdate(user *Model.User, id uint) (*Model.User, error) {
	var userModel Model.User
	result := r.db.First(&userModel, id)
	if result.Error != nil {
		return nil, result.Error
	}
	userModel.SetUsername(user.GetUsername())
	result = r.db.Save(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (r *UserRepository) UserCreate(user *Model.User) (*Model.User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepository) FindById(id uint) (*Model.User, error) {
	var user Model.User
	err := r.db.Preload("Posters").Preload("Posters.Images").Preload("Posters.Addresses").Preload("Posters.Tags").
		Preload("MarkedPosters").Preload("MarkedPosters.Poster").Preload("MarkedPosters.Poster.Images").Preload("MarkedPosters.Poster.Addresses").Preload("MarkedPosters.Poster.Tags").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() (*[]Model.User, error) {
	var users []Model.User
	err := r.db.Preload("Posters").Preload("Posters.Images").Preload("Posters.Addresses").Preload("Posters.Tags").
		Preload("MarkedPosters").Preload("MarkedPosters.Poster").Preload("MarkedPosters.Poster.Images").Preload("MarkedPosters.Poster.Addresses").Preload("MarkedPosters.Poster.Tags").
		Find(&users).Error
	if err != nil {
		return nil, errors.New("error while getting all users")
	}
	return &users, nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	var user Model.User
	if err := r.db.Preload("Posters").Preload("OwnConversations").Preload("MemberConversations").
		Preload("MarkedPosters").Preload("Payments").First(&user, id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&user.Posters, "user_id = ?", id).Error; err != nil {
		return err
	}
	// delete images and addresses
	for _, poster := range user.Posters {
		if err := r.db.Delete(&poster.Images, "poster_id = ?", poster.ID).Error; err != nil {
			return err
		}
		if err := r.db.Delete(&poster.Addresses, "poster_id = ?", poster.ID).Error; err != nil {
			return err
		}
	}
	if err := r.db.Delete(&user.MarkedPosters, "user_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&user.Payments, "user_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&user.OwnConversations, "owner_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&user.MemberConversations, "member_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateWallet(id uint, amount float64) (*Model.User, error) {
	var user Model.User
	result := r.db.Model(&user).Where("id = ?", id).Update("wallet", gorm.Expr("wallet + ?", amount))
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepository) GetAmount(id uint) (float64, error) {
	var user Model.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.GetWallet(), nil
}
