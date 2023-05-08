package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAdminById(id int) (Model.Admin, error) {
	var admin Model.Admin
	result := r.db.First(&admin, id)
	if result.Error != nil {
		return Model.Admin{}, result.Error
	}

	return admin, nil
}

func (r *AdminRepository) CreateAdmin(admin Model.Admin) (Model.Admin, error) {
	result := r.db.Create(&admin)
	if result.Error != nil {
		return Model.Admin{}, result.Error
	}

	return admin, nil
}

func (r *AdminRepository) GetAdminByUsername(username string) (Model.Admin, error) {
	var admin Model.Admin
	result := r.db.Where("username = ?", username).First(&admin)
	if result.Error != nil {
		return Model.Admin{}, result.Error
	}

	return admin, nil
}
