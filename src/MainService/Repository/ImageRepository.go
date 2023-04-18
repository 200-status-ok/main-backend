package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) CreateImage(imageUrl []string, posterID uint) ([]Model.Image, error) {
	var imageModels []Model.Image
	for _, url := range imageUrl {
		var image Model.Image
		image.SetPosterID(posterID)
		image.SetUrl(url)
		imageModels = append(imageModels, image)
	}
	result := r.db.Create(&imageModels)
	if result.Error != nil {
		return []Model.Image{}, result.Error
	}

	return imageModels, nil
}

func (r *ImageRepository) DeleteImageById(id int) error {
	var imageModel Model.Image
	result := r.db.Delete(&imageModel, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ImageRepository) DeleteImageByPosterId(posterId uint) error {
	var imageModel Model.Image
	result := r.db.Delete(&imageModel, "poster_id = ?", posterId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ImageRepository) UpdateImage(imageUrl []string, posterID uint) ([]Model.Image, error) {
	var imageModel []Model.Image
	result := r.db.First(&imageModel, "poster_id = ?", posterID)
	if result.Error != nil {
		return []Model.Image{}, result.Error
	}
	for i, url := range imageUrl {
		imageModel[i].SetUrl(url)
	}
	result = r.db.Save(&imageModel)
	if result.Error != nil {
		return []Model.Image{}, result.Error
	}

	return imageModel, nil
}
