package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository { //todo modar singleton
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

func (r *CategoryRepository) GetTagByName(name string) (Model.Tag, error) {
	var category Model.Tag
	result := r.db.Where("name = ?", name).First(&category)
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

func (r *CategoryRepository) UpdateTag(id uint, tag Model.Tag) error {
	var updatedTagModel Model2.Tag

	if tag.Name != "" {
		updatedTagModel.Name = tag.Name
	}

	if tag.State != "" {
		updatedTagModel.State = tag.State
	}

	result := r.db.Model(&Model2.Tag{}).Where("id = ?", id).Updates(updatedTagModel)

	if result.Error != nil {
		return result.Error
	}

	return nil
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

func (r *CategoryRepository) GetTags(state string) ([]Model2.Tag, error) {
	var categories []Model2.Tag

	var result *gorm.DB

	if state != "" && state != "all" {
		result = r.db.Where("state = ?", state).Find(&categories)
	} else {
		result = r.db.Find(&categories)
	}

	if result.Error != nil {
		return []Model2.Tag{}, result.Error
	}

	return categories, nil
}
