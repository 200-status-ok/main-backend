package Repository

import (
	"github.com/403-access-denied/main-backend/src/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/DTO"
	"github.com/403-access-denied/main-backend/src/Model"
	"gorm.io/gorm"
)

type PosterRepository struct {
	db *gorm.DB
}

func NewPosterRepository(db *gorm.DB) *PosterRepository {
	return &PosterRepository{db: db}
}

func (r *PosterRepository) GetAllPosters(limit, offset int, sort, sortBy string) ([]Model.Poster, error) {
	var posters []Model.Poster
	result := r.db.Preload("Address").Preload("Images").Preload("Categories").Preload("User").
		Limit(limit).Offset(offset).Order(sortBy + " " + sort).Find(&posters)
	DBConfiguration.CloseDB()
	if result.Error != nil {
		return nil, result.Error
	}
	return posters, nil
}

func (r *PosterRepository) GetPosterById(id int) (Model.Poster, error) {
	var poster Model.Poster
	result := r.db.Preload("Address").Preload("Images").Preload("Categories").Preload("User").
		First(&poster, "id = ?", id)
	DBConfiguration.CloseDB()
	if result.Error != nil {
		return Model.Poster{}, result.Error
	}
	return poster, nil
}

func (r *PosterRepository) DeletePosterById(id int) error {
	var poster Model.Poster
	result := r.db.Find(&poster, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if err := r.db.Delete(&poster.Address, "poster_id = ?", id).Error; err != nil {
		return err
	}
	result = r.db.Where("poster_id = ?", id).Delete(&Model.Image{})
	if result.Error != nil {
		return result.Error
	}
	_ = r.db.Model(&poster).Association("Categories").Clear()
	if err := r.db.Delete(&poster, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PosterRepository) CreatePoster(poster DTO.PosterDTO, address DTO.AddressDTO, imageUrl []string, categories []int) (
	Model.Poster, error) {
	var posterModel Model.Poster
	var categoriesModel []Model.Category
	posterModel.SetTitle(poster.Title)
	posterModel.SetDescription(poster.Description)
	posterModel.SetUserID(poster.UserID)
	posterModel.SetStatus(poster.Status)
	posterModel.SetUserPhone(poster.UserPhone)
	posterModel.SetTelegramID(poster.TelID)
	posterModel.SetHasAlert(poster.Alert)
	posterModel.SetAward(poster.Award)

	for _, category := range categories {
		categoryModel, err := NewCategoryRepository(r.db).GetCategoryById(category)
		if err != nil {
			return Model.Poster{}, err
		}
		categoriesModel = append(categoriesModel, categoryModel)
	}
	posterModel.SetCategories(categoriesModel)

	result := r.db.Create(&posterModel)
	if result.Error != nil {
		return Model.Poster{}, result.Error
	}

	posterID := posterModel.GetID()
	newAddress, err := NewAddressRepository(r.db).CreateAddress(address, posterID)
	if err != nil {
		return Model.Poster{}, err
	}
	posterModel.SetAddress(newAddress)
	newImages, err := NewImageRepository(r.db).CreateImage(imageUrl, posterID)
	if err != nil {
		return Model.Poster{}, err
	}
	posterModel.SetImages(newImages)

	return posterModel, nil
}

func (r *PosterRepository) UpdatePoster(id int, poster DTO.PosterDTO, address DTO.AddressDTO, imageUrl []string, categories []int) (
	Model.Poster, error) {
	var posterModel Model.Poster
	result := r.db.Preload("Address").Preload("Images").Preload("Categories").Preload("User").
		First(&posterModel, "id = ?", id)
	if result.Error != nil {
		return Model.Poster{}, result.Error
	}
	var categoriesModel []Model.Category
	posterModel.SetTitle(poster.Title)
	posterModel.SetDescription(poster.Description)
	posterModel.SetUserID(poster.UserID)
	posterModel.SetStatus(poster.Status)
	posterModel.SetUserPhone(poster.UserPhone)
	posterModel.SetTelegramID(poster.TelID)
	posterModel.SetHasAlert(poster.Alert)
	posterModel.SetAward(poster.Award)

	for _, category := range categories {
		categoryModel, err := NewCategoryRepository(r.db).GetCategoryById(category)
		if err != nil {
			return Model.Poster{}, err
		}
		categoriesModel = append(categoriesModel, categoryModel)
	}
	posterModel.SetCategories(categoriesModel)

	result = r.db.Save(&posterModel)
	if result.Error != nil {
		return Model.Poster{}, result.Error
	}

	posterID := posterModel.GetID()
	updatedAddress, err := NewAddressRepository(r.db).UpdateAddress(address, posterID)
	if err != nil {
		return Model.Poster{}, err
	}
	posterModel.SetAddress(updatedAddress)
	updatedImage, err := NewImageRepository(r.db).UpdateImage(imageUrl, posterID)
	if err != nil {
		return Model.Poster{}, err
	}
	posterModel.SetImages(updatedImage)

	return posterModel, nil
}
