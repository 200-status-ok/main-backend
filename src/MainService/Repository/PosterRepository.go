package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	DTO2 "github.com/403-access-denied/main-backend/src/MainService/DTO"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"gorm.io/gorm"
	"sort"
	"strings"
)

type PosterRepository struct {
	db *gorm.DB
}

func NewPosterRepository(db *gorm.DB) *PosterRepository {
	return &PosterRepository{db: db}
}

type ScoredPoster struct {
	poster Model2.Poster
	score  int
}

func (r *PosterRepository) GetAllPosters(limit, offset int, sortType, sortBy string, filterObject DTO2.FilterObject) ([]Model2.Poster, error) {
	var posters []Model2.Poster

	var result *gorm.DB
	if filterObject.SearchPhrase != "" || filterObject.TagIds != nil {
		result = r.db.Preload("Addresses").Preload("Images").Preload("Categories").Preload("User").
			Order(sortBy + " " + sortType)
	} else {
		result = r.db.Preload("Addresses").Preload("Images").Preload("Categories").Preload("User").
			Limit(limit).Offset(offset).Order(sortBy + " " + sortType)
	}

	if filterObject.Status != "" && filterObject.Status != "both" {
		result = result.Where("status = ?", filterObject.Status)
	}

	if filterObject.OnlyRewards == true {
		result = result.Where("award > 0")
	}

	if filterObject.TimeStart != 0 && filterObject.TimeEnd != 0 {
		result = result.Where("extract(epoch from posters.created_at) BETWEEN ? AND ?", filterObject.TimeStart, filterObject.TimeEnd)
	}

	if filterObject.Lat != 0 && filterObject.Lon != 0 {
		result = result.Select("posters.*").Joins("LEFT JOIN addresses ON posters.id = addresses.poster_id AND addresses.deleted_at IS NULL").
			Where("Addresses.latitude BETWEEN ? AND ? AND Addresses.longitude BETWEEN ? AND ? AND posters.deleted_at IS NULL",
				filterObject.Lat-Utils.BaseLocationRadarRadius, filterObject.Lat+Utils.BaseLocationRadarRadius,
				filterObject.Lon-Utils.BaseLocationRadarRadius, filterObject.Lon+Utils.BaseLocationRadarRadius)
	}

	result.Find(&posters)

	if filterObject.SearchPhrase != "" || filterObject.TagIds != nil {

		var validPosters []ScoredPoster

		for _, poster := range posters {
			searchScore := 0

			if filterObject.SearchPhrase != "" {

				if strings.Contains(filterObject.SearchPhrase, poster.Title) || strings.Contains(poster.Title, filterObject.SearchPhrase) {
					searchScore += 3
				}

				if strings.Contains(poster.Description, filterObject.SearchPhrase) {
					searchScore += 1
				}

				for _, category := range poster.Categories {
					if strings.Contains(category.Name, filterObject.SearchPhrase) || strings.Contains(filterObject.SearchPhrase, category.Name) {
						searchScore += 1
					}
				}
			}

			if filterObject.TagIds != nil {
				for _, tagId := range filterObject.TagIds {
					for _, category := range poster.Categories {
						if tagId == int(category.ID) {
							searchScore += 1
						}
					}
				}
			}
			if searchScore != 0 {
				validPosters = append(validPosters, ScoredPoster{
					poster: poster,
					score:  searchScore,
				})
			}
		}

		sort.Slice(validPosters, func(i, j int) bool {
			return validPosters[i].score < validPosters[j].score
		})

		var sortedPosters []Model2.Poster
		for _, s := range validPosters {
			sortedPosters = append(sortedPosters, s.poster)
		}

		if offset+limit > len(sortedPosters) {
			return sortedPosters[offset:], nil
		} else {
			return sortedPosters[offset : offset+limit], nil
		}
	}

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
	posterModel.HasChat = poster.Chat

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

	if posterModel.HasChat {
		err := NewChatRepository(r.db).CreateChatRoom(posterID, poster.UserID)
		if err != nil {
			return Model2.Poster{}, err
		}
	}

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

func (r *PosterRepository) CreatePosterReport(posterID uint, issuerID uint, reportType string, description string) error {

	var reportModel = Model2.PosterReport{
		PosterID:    posterID,
		IssuerID:    issuerID,
		ReportType:  reportType,
		Description: description,
		Status:      "open",
	}

	result := r.db.Create(&reportModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PosterRepository) GetAllPosterReports(limit, offset int, status string) ([]Model2.PosterReport, error) {
	var posterReports []Model2.PosterReport

	var result = r.db.Preload("Poster").Preload("Issuer").Preload("Poster.User").Where("deleted_at IS NULL").Limit(limit).Offset(offset).Order("created_at DESC")

	if status != "both" {
		result = result.Where("status = ?", status)
	}

	result.Find(&posterReports)

	DBConfiguration.CloseDB()
	if result.Error != nil {
		return nil, result.Error
	}
	return posterReports, nil
}

func (r *PosterRepository) GetPosterReportById(id int) (Model2.PosterReport, error) {
	var posterReport Model2.PosterReport

	result := r.db.Preload("Poster").Preload("Issuer").First(&posterReport, "id = ?", id)
	DBConfiguration.CloseDB()

	if result.Error != nil {
		return Model2.PosterReport{}, result.Error
	}

	return posterReport, nil
}

func (r *PosterRepository) UpdatePosterReport(id, posterID, issuerID uint, reportType, description, status string) error {

	var updatedPosterReportModel Model2.PosterReport

	if posterID != 0 {
		updatedPosterReportModel.PosterID = posterID
	}

	if issuerID != 0 {
		updatedPosterReportModel.IssuerID = issuerID
	}

	if reportType != "" {
		updatedPosterReportModel.ReportType = reportType
	}

	if description != "" {
		updatedPosterReportModel.Description = description
	}

	if status != "" {
		updatedPosterReportModel.Status = status
	}

	result := r.db.Model(&Model2.PosterReport{}).Where("id = ?", id).Updates(updatedPosterReportModel)

	if result.Error != nil {
		return result.Error
	}

	DBConfiguration.CloseDB()

	return nil
}
