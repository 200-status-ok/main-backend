package Repository

import (
	"errors"
	"fmt"
	DTO2 "github.com/403-access-denied/main-backend/src/MainService/DTO"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		result = r.db.Preload("Addresses").Preload("Images").Preload("Tags", "state = ?", "accepted").
			Order(sortBy + " " + sortType)
	} else {
		result = r.db.Preload("Addresses").Preload("Images").Preload("Tags", "state = ?", "accepted").
			Limit(limit).Offset(offset).Order(sortBy + " " + sortType)
	}

	if filterObject.State != "" && filterObject.State != "all" {
		result = result.Where("state = ?", filterObject.State)
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

				for _, category := range poster.Tags {
					if strings.Contains(category.Name, filterObject.SearchPhrase) || strings.Contains(filterObject.SearchPhrase, category.Name) {
						searchScore += 1
					}
				}
			}

			if filterObject.TagIds != nil {
				for _, tagId := range filterObject.TagIds {
					for _, category := range poster.Tags {
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

	if result.Error != nil {
		return nil, result.Error
	}
	return posters, nil
}

func (r *PosterRepository) GetPosterById(id int) (Model2.Poster, error) {
	var poster Model2.Poster
	result := r.db.Preload("Addresses").Preload("Images").Preload("Tags", "state = ?", "accepted").
		First(&poster, "id = ?", id)

	if result.Error != nil {
		return Model2.Poster{}, result.Error
	}
	return poster, nil
}

func (r *PosterRepository) DeletePosterById(id uint, userId uint) error {
	var poster Model2.Poster

	findResult := r.db.First(&poster, "id = ? AND user_id = ?", id, userId)
	if findResult.Error != nil {
		return findResult.Error
	}

	deleteResult := r.db.Select(clause.Associations).Delete(&poster)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	return nil

}

func (r *PosterRepository) CreatePoster(poster DTO2.CreatePosterDTO, addresses []DTO2.CreateAddressDTO, imageUrls []string, tagNames []string) (
	Model2.Poster, error) {
	var posterModel Model2.Poster
	posterModel.SetTitle(poster.Title)
	posterModel.SetDescription(poster.Description)
	posterModel.SetUserID(poster.UserID)
	posterModel.SetStatus(poster.Status)
	posterModel.SetUserPhone(poster.UserPhone)
	posterModel.SetTelegramID(poster.TelID)
	posterModel.SetHasAlert(poster.Alert)
	posterModel.SetAward(poster.Award)
	posterModel.SetHasChat(poster.Chat)
	posterModel.SetState("pending")

	var newTags []Model2.Tag
	tagRepository := NewCategoryRepository(r.db)

	for _, tagName := range tagNames {
		name := strings.ToLower(strings.Trim(tagName, " "))
		tagModel, getErr := tagRepository.GetTagByName(name)
		if getErr != nil {
			newTag, creErr := tagRepository.CreateCategory(Model2.Tag{
				Name:  name,
				State: "pending",
			})
			if creErr != nil {
				continue
			}
			newTags = append(newTags, newTag)
		} else {
			newTags = append(newTags, tagModel)
		}
	}
	_ = r.db.Model(&posterModel).Association("Tags").Append(newTags)

	var newAddresses []Model2.Address
	for _, address := range addresses {
		newAddress := Model2.Address{
			Province:      address.Province,
			City:          address.City,
			AddressDetail: address.AddressDetail,
			Latitude:      address.Latitude,
			Longitude:     address.Longitude,
		}
		newAddresses = append(newAddresses, newAddress)
	}
	_ = r.db.Model(&posterModel).Association("Addresses").Append(newAddresses)

	var newImages []Model2.Image
	for _, url := range imageUrls {
		image := Model2.Image{
			Url: url,
		}
		newImages = append(newImages, image)
	}
	_ = r.db.Model(&posterModel).Association("Images").Append(newImages)

	result := r.db.Create(&posterModel)
	if result.Error != nil {
		return Model2.Poster{}, result.Error
	}

	return posterModel, nil
}

func (r *PosterRepository) UpdatePoster(id int, poster DTO2.UpdatePosterDTO, addresses []DTO2.UpdateAddressDTO) error {

	var updatedPosterModel Model2.Poster
	updatedPosterModel.SetID(uint(id))

	if poster.Title != "" {
		updatedPosterModel.Title = poster.Title
	}

	if poster.Description != "" {
		updatedPosterModel.Description = poster.Description
	}

	if poster.Status != "" {
		updatedPosterModel.SetStatus(poster.Status)
	}

	if poster.TelID != "" {
		updatedPosterModel.TelegramID = poster.TelID
	}

	if poster.UserPhone != "" {
		updatedPosterModel.UserPhone = poster.UserPhone
	}

	fmt.Println("modar alert: ", poster.Alert)
	if poster.Alert != "" { //todo modar fix alert
		if poster.Alert == "true" {
			updatedPosterModel.HasAlert = true
		} else if poster.Alert == "false" {
			fmt.Println("modar alert is false")
			updatedPosterModel.HasAlert = false
		}
	}

	fmt.Println("modar Chat: ", poster.Chat)
	if poster.Chat != "" { //todo modar fix chat
		if poster.Chat == "true" {
			updatedPosterModel.HasChat = true
		} else if poster.Chat == "false" {
			updatedPosterModel.HasChat = false
		}
	}

	if poster.Award != 0 {
		updatedPosterModel.Award = poster.Award
	}

	if poster.UserID != 0 {
		updatedPosterModel.UserID = poster.UserID
	}

	if poster.State != "" {
		updatedPosterModel.State = poster.State
	}

	if len(poster.ImgUrls) != 0 {
		var updatedImages []Model2.Image
		for _, url := range poster.ImgUrls {
			updatedImages = append(updatedImages, Model2.Image{
				Url:      url,
				PosterID: uint(id),
			})
		}
		r.db.Where("poster_id = ?", id).Delete(&Model2.Image{})
		_ = r.db.Model(&updatedPosterModel).Association("Images").Append(updatedImages)
	}

	//todo modar add tag and address support
	//if len(poster.TagIds) != 0 {
	//	var updatedTags []Model2.Tag
	//	for _, tagId := range poster.TagIds {
	//		updatedTags = append(updatedTags, Model2.Tag{
	//			ID: tagId,
	//		})
	//	}
	//	_ = r.db.Model(&updatedPosterModel).Association("Tags").Replace(updatedTags)
	//}
	//
	//if len(addresses) != 0 {
	//	var updatedAddresses []Model2.Address
	//	for _, address := range addresses {
	//		updatedAddresses = append(updatedAddresses, Model2.Address{
	//			Province:      address.Province,
	//			City:          address.City,
	//			AddressDetail: address.AddressDetail,
	//			Latitude:      address.Latitude,
	//			Longitude:     address.Longitude,
	//		})
	//	}
	//	_ = r.db.Model(&updatedPosterModel).Association("Addresses").Replace(updatedAddresses)
	//}

	fmt.Println("modar updatedPosterModel: ", updatedPosterModel)

	//if len(poster.TagIds) != 0 {
	//	var updatedTags []Model2.Tag
	//	for _, tagId := range poster.TagIds {
	//		updatedTags = append(updatedTags, Model2.Tag{
	//			Url: url,
	//		})
	//	}
	//	_ = r.db.Model(&updatedPosterModel).Association("Images").Replace(updatedTags)
	//}

	fmt.Println("modar updateStatus: ", updatedPosterModel.Status)

	result := r.db.Where("id = ?", id).Updates(updatedPosterModel)

	if result.RowsAffected == 0 {
		return errors.New("poster not found")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
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

	var result = r.db.Preload("Poster").Preload("Issuer").Preload("Poster.User").Limit(limit).Offset(offset).Order("created_at DESC")

	if status != "both" {
		result = result.Where("status = ?", status)
	}

	result.Find(&posterReports)

	if result.Error != nil {
		return nil, result.Error
	}
	return posterReports, nil
}

func (r *PosterRepository) GetPosterReportById(id int) (Model2.PosterReport, error) {
	var posterReport Model2.PosterReport

	result := r.db.Preload("Poster").Preload("Issuer").First(&posterReport, "id = ?", id)

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

	return nil
}

func (r *PosterRepository) UpdatePosterState(id uint, state string) error {

	var updatedPosterReportModel Model2.Poster

	updatedPosterReportModel.State = state

	fmt.Println("modar id is: ", id)
	fmt.Println("modar updatedPosterReportModel is: ", updatedPosterReportModel)

	result := r.db.Model(&Model2.Poster{}).Where("id = ?", id).Updates(updatedPosterReportModel)

	if result.Error != nil {
		return result.Error
	}

	//DBConfiguration.CloseDB()

	return nil
}
