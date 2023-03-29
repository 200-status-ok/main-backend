package Repository

import (
	"github.com/403-access-denied/main-backend/src/Model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategoryById(id int) (Model.Category, error) {
	var category Model.Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		return Model.Category{}, result.Error
	}

	return category, nil
}
