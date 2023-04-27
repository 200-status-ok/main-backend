package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
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

func (r *CategoryRepository) CreateCategory(category Model.Category) (Model.Category, error) {
	result := r.db.Create(&category)
	if result.Error != nil {
		return Model.Category{}, result.Error
	}

	return category, nil
}

func (r *CategoryRepository) UpdateCategory(id uint, category Model.Category) (Model.Category, error) {
	//result := r.db.Model(&category).Where("id = ?", id).Updates(category)
	//if result.Error != nil {
	//	return Model.Category{}, result.Error
	//}
	//
	//return category, nil

	var categoryModel Model.Category
	result := r.db.First(&categoryModel, id)
	if result.Error != nil {
		return Model.Category{}, result.Error
	}
	categoryModel.SetName(category.GetName())
	result = r.db.Save(&categoryModel)
	if result.Error != nil {
		return Model.Category{}, result.Error
	}

	return categoryModel, nil
}

func (r *CategoryRepository) DeleteCategory(id uint) error {
	var categoryModel Model.Category
	if err := r.db.Where("id = ?", id).First(&categoryModel).Error; err != nil {
		return err
	}
	result := r.db.Delete(&categoryModel, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CategoryRepository) GetCategories() ([]Model2.Category, error) {
	var categories []Model2.Category
	result := r.db.Find(&categories)
	if result.Error != nil {
		return []Model2.Category{}, result.Error
	}

	return categories, nil
}
