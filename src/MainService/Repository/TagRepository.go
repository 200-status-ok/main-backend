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

func (r *CategoryRepository) GetCategoryById(id int) (Model.Tag, error) {
	var category Model.Tag
	result := r.db.First(&category, id)
	if result.Error != nil {
		return Model.Tag{}, result.Error
	}

	return category, nil
}

func (r *CategoryRepository) CreateCategory(category Model.Tag) (Model.Tag, error) {
	result := r.db.Create(&category)
	if result.Error != nil {
		return Model.Tag{}, result.Error
	}

	return category, nil
}

func (r *CategoryRepository) UpdateCategory(id uint, category Model.Tag) (Model.Tag, error) {
	//result := r.db.Model(&category).Where("id = ?", id).Updates(category)
	//if result.Error != nil {
	//	return Model.Tag{}, result.Error
	//}
	//
	//return category, nil

	var categoryModel Model.Tag
	result := r.db.First(&categoryModel, id)
	if result.Error != nil {
		return Model.Tag{}, result.Error
	}
	categoryModel.SetName(category.GetName())
	result = r.db.Save(&categoryModel)
	if result.Error != nil {
		return Model.Tag{}, result.Error
	}

	return categoryModel, nil
}

func (r *CategoryRepository) DeleteCategory(id uint) error {
	var categoryModel Model.Tag
	result := r.db.Find(&categoryModel, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	result = r.db.Where("poster_id = ?", id).Delete(&Model2.Image{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CategoryRepository) GetCategories() ([]Model2.Tag, error) {
	var categories []Model2.Tag
	result := r.db.Find(&categories)
	if result.Error != nil {
		return []Model2.Tag{}, result.Error
	}

	return categories, nil
}
