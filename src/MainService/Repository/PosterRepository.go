package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	DTO2 "github.com/403-access-denied/main-backend/src/MainService/DTO"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"gorm.io/gorm"
)

type PosterRepository struct {
	db *gorm.DB
}

func NewPosterRepository(db *gorm.DB) *PosterRepository {
	return &PosterRepository{db: db}
}

func (r *PosterRepository) GetAllPosters(limit, offset int, sort, sortBy string, filterObject DTO2.FilterObject) ([]Model2.Poster, error) {
	var posters []Model2.Poster

	result := r.db.Preload("Addresses").Preload("Images").Preload("Categories").Preload("User").
		Limit(limit).Offset(offset).Order(sortBy + " " + sort)

	if filterObject.Status != "" && filterObject.Status != "both" { // todo use enum maybe?
		result = result.Where("status = ?", filterObject.Status)
	}

	if filterObject.SearchPhrase != "" {
		result = result.Where("title LIKE ?", "%"+filterObject.SearchPhrase+"%")
	}

	if filterObject.Lat != 0 && filterObject.Lon != 0 {
		result = result.Where("Lat BETWEEN ? AND ?", filterObject.Lat-Utils.BaseLocationRadarRadius, filterObject.Lat+Utils.BaseLocationRadarRadius).
			Where("Lon BETWEEN ? AND ?", filterObject.Lon-Utils.BaseLocationRadarRadius, filterObject.Lon+Utils.BaseLocationRadarRadius) // todo change logic maybe
	}

	if filterObject.TimeStart != 0 && filterObject.TimeEnd != 0 {
		result = result.Where("created_at BETWEEN ? AND ?", filterObject.TimeStart, filterObject.TimeEnd)
	}

	if filterObject.OnlyRewards {
		result = result.Where("only_awards = ?", filterObject.OnlyRewards)
	}

	result.Find(&posters)
	DBConfiguration.CloseDB()
	if result.Error != nil {
		return nil, result.Error
	}
	return posters, nil
}

func (r *PosterRepository) GetPosterById(id int) (Model2.Poster, error) {
	var poster Model2.Poster
	result := r.db.Preload("Addresses").Preload("Images").Preload("Categories").Preload("User").
		First(&poster, "id = ?", id)
	DBConfiguration.CloseDB()
	if result.Error != nil {
		return Model2.Poster{}, result.Error
	}
	return poster, nil
}

func (r *PosterRepository) DeletePosterById(id int) error {
	var poster Model2.Poster
	result := r.db.Find(&poster, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if err := r.db.Delete(&poster.Addresses, "poster_id = ?", id).Error; err != nil {
		return err
	}
	result = r.db.Where("poster_id = ?", id).Delete(&Model2.Image{})
	if result.Error != nil {
		return result.Error
	}
	_ = r.db.Model(&poster).Association("Categories").Clear()
	if err := r.db.Delete(&poster, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PosterRepository) CreatePoster(poster DTO2.PosterDTO, addresses []DTO2.AddressDTO, imageUrl []string, categories []int) (
	Model2.Poster, error) {
	var posterModel Model2.Poster
	var categoriesModel []Model2.Category
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
			return Model2.Poster{}, err
		}
		categoriesModel = append(categoriesModel, categoryModel)
	}
	posterModel.SetCategories(categoriesModel)

	result := r.db.Create(&posterModel)
	if result.Error != nil {
		return Model2.Poster{}, result.Error
	}

	posterID := posterModel.GetID()
	newAddress, err := NewAddressRepository(r.db).CreateAddress(addresses, posterID)
	if err != nil {
		return Model2.Poster{}, err
	}
	posterModel.SetAddress(newAddress)
	newImages, err := NewImageRepository(r.db).CreateImage(imageUrl, posterID)
	if err != nil {
		return Model2.Poster{}, err
	}
	posterModel.SetImages(newImages)

	return posterModel, nil
}

func (r *PosterRepository) UpdatePoster(id int, poster DTO2.PosterDTO, addresses []DTO2.AddressDTO, imageUrl []string, categories []int) (
	Model2.Poster, error) {
	var posterModel Model2.Poster
	result := r.db.Preload("Addresses").Preload("Images").Preload("Categories").Preload("User").
		First(&posterModel, "id = ?", id)
	if result.Error != nil {
		return Model2.Poster{}, result.Error
	}
	var categoriesModel []Model2.Category
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
			return Model2.Poster{}, err
		}
		categoriesModel = append(categoriesModel, categoryModel)
	}
	posterModel.SetCategories(categoriesModel)

	result = r.db.Save(&posterModel)
	if result.Error != nil {
		return Model2.Poster{}, result.Error
	}

	posterID := posterModel.GetID()
	updatedAddress, err := NewAddressRepository(r.db).UpdateAddress(addresses, posterID)
	if err != nil {
		return Model2.Poster{}, err
	}
	posterModel.SetAddress(updatedAddress)
	updatedImage, err := NewImageRepository(r.db).UpdateImage(imageUrl, posterID)
	if err != nil {
		return Model2.Poster{}, err
	}
	posterModel.SetImages(updatedImage)

	return posterModel, nil
}
