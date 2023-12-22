package Repository

import (
	"errors"
	"fmt"
	Model2 "github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	DTO2 "github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
	"strings"
	"time"
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

func (r *PosterRepository) GetAllESPosters() ([]*DTO2.ESPosterDTO, error) {
	var posters []Model2.Poster
	result := r.db.Preload("Addresses").Preload("Images").Preload("Tags").Find(&posters).Error
	if result != nil {
		return []*DTO2.ESPosterDTO{}, result
	}

	var esPosters []*DTO2.ESPosterDTO
	for _, poster := range posters {
		var esTag []DTO2.ESTagDTO
		for _, tag := range poster.Tags {
			esTag = append(esTag, DTO2.ESTagDTO{
				ID:    tag.ID,
				Name:  tag.Name,
				State: tag.State,
			})
		}
		var esAddress []DTO2.ESAddressDTO
		for _, address := range poster.Addresses {
			var location DTO2.Location
			location.Latitude = address.Latitude
			location.Longitude = address.Longitude
			esAddress = append(esAddress, DTO2.ESAddressDTO{
				Province:      address.Province,
				City:          address.City,
				AddressDetail: address.AddressDetail,
				Location:      location,
			})
		}
		var imageUrls []string
		for _, image := range poster.Images {
			imageUrls = append(imageUrls, image.Url)
		}
		esPosters = append(esPosters, &DTO2.ESPosterDTO{
			ID:          poster.ID,
			Title:       poster.Title,
			Description: poster.Description,
			Status:      string(poster.Status),
			TelID:       poster.TelegramID,
			UserPhone:   poster.UserPhone,
			Alert:       poster.HasAlert,
			Chat:        poster.HasChat,
			Award:       poster.Award,
			UserID:      poster.UserID,
			State:       poster.State,
			SpecialType: poster.SpecialType,
			CreatedAt:   poster.CreatedAt.Unix(),
			UpdatedAt:   poster.UpdatedAt.Unix(),
			Addresses:   esAddress,
			Images:      imageUrls,
			Tags:        esTag,
		})
	}

	return esPosters, nil
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

	if filterObject.SpecialType != "" && filterObject.SpecialType != "all" {
		result = result.Where("special_type = ?", filterObject.SpecialType)
	}
	if filterObject.State != "" && filterObject.State != "all" {
		result = result.Where("state = ?", filterObject.State)
	}

	if filterObject.Status != "" && filterObject.Status != "both" {
		result = result.Where("status = ?", filterObject.Status)
	}

	if filterObject.OnlyAwards == true {
		result = result.Where("award > 0")
	}

	if filterObject.TimeStart != 0 && filterObject.TimeEnd != 0 {
		result = result.Where("extract(epoch from posters.created_at) BETWEEN ? AND ?", filterObject.TimeStart, filterObject.TimeEnd)
	}

	if filterObject.Lat != 0 && filterObject.Lon != 0 {
		result = result.Select("posters.*").Joins("LEFT JOIN addresses ON posters.id = addresses.poster_id AND addresses.deleted_at IS NULL").
			Where("Addresses.latitude BETWEEN ? AND ? AND Addresses.longitude BETWEEN ? AND ? AND posters.deleted_at IS NULL",
				filterObject.Lat-utils.BaseLocationRadarRadius, filterObject.Lat+utils.BaseLocationRadarRadius,
				filterObject.Lon-utils.BaseLocationRadarRadius, filterObject.Lon+utils.BaseLocationRadarRadius)
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
	esPosterCli := ElasticSearch.NewPosterES(elasticsearch.GetElastic())

	findResult := r.db.First(&poster, "id = ? AND user_id = ?", id, userId)
	if findResult.Error != nil {
		return findResult.Error
	}

	deleteResult := r.db.Select(clause.Associations).Delete(&poster)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	err := esPosterCli.DeletePoster(int(poster.ID))
	if err != nil {
		return err
	}

	return nil

}

func (r *PosterRepository) CreatePoster(userID uint64, poster DTO2.CreatePosterDTO, addresses []DTO2.CreateAddressDTO,
	imageUrls []string, tagNames []string, special string) (
	Model2.Poster, error) {
	var posterModel Model2.Poster
	posterModel.SetTitle(poster.Title)
	posterModel.SetDescription(poster.Description)
	posterModel.SetUserID(uint(userID))
	posterModel.SetStatus(poster.Status)
	posterModel.SetUserPhone(poster.UserPhone)
	posterModel.SetTelegramID(poster.TelID)
	posterModel.SetHasAlert(poster.Alert)
	posterModel.SetAward(poster.Award)
	posterModel.SetHasChat(poster.Chat)
	posterModel.SetState("pending")
	posterModel.SetSpecialType(special)

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

	esPosterCli := ElasticSearch.NewPosterES(elasticsearch.GetElastic())

	var esPoster DTO2.ESPosterDTO
	esPoster.ID = posterModel.ID
	esPoster.Title = posterModel.Title
	esPoster.Description = posterModel.Description
	esPoster.Status = string(posterModel.Status)
	esPoster.UserID = posterModel.UserID
	esPoster.UserPhone = posterModel.UserPhone
	esPoster.TelID = posterModel.TelegramID
	esPoster.Alert = posterModel.HasAlert
	esPoster.Award = posterModel.Award
	esPoster.Chat = posterModel.HasChat
	esPoster.State = posterModel.State
	esPoster.SpecialType = posterModel.SpecialType

	var esTags []DTO2.ESTagDTO
	for _, v := range newTags {
		esTags = append(esTags, DTO2.ESTagDTO{
			ID:    v.ID,
			Name:  v.Name,
			State: v.State,
		})
	}
	esPoster.Tags = esTags

	var esAddresses []DTO2.ESAddressDTO
	for _, v := range addresses {
		var location DTO2.Location
		location.Latitude = v.Latitude
		location.Longitude = v.Longitude

		esAddresses = append(esAddresses, DTO2.ESAddressDTO{
			Province:      v.Province,
			City:          v.City,
			AddressDetail: v.AddressDetail,
			Location:      location,
		})
	}

	esPoster.Addresses = esAddresses
	esPoster.CreatedAt = posterModel.CreatedAt.Unix()
	esPoster.UpdatedAt = posterModel.UpdatedAt.Unix()
	esPoster.Images = imageUrls

	err := esPosterCli.InsertPoster(&esPoster)
	if err != nil {
		return posterModel, err
	}

	return posterModel, nil
}

func (r *PosterRepository) UpdatePoster(id int, role string, poster DTO2.UpdatePosterDTO, addresses []DTO2.UpdateAddressDTO) error {

	var updatedPosterModel Model2.Poster
	updatedPosterModel.SetID(uint(id))

	esPosterCli := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	updateFields := make(map[string]interface{})

	if poster.Title != "" {
		updatedPosterModel.Title = poster.Title
		updateFields["title"] = poster.Title
	}

	if poster.Description != "" {
		updatedPosterModel.Description = poster.Description
		updateFields["description"] = poster.Description
	}

	if poster.Status != "" {
		updatedPosterModel.SetStatus(poster.Status)
		updateFields["status"] = poster.Status
	}

	if poster.TelID != "" {
		updatedPosterModel.TelegramID = poster.TelID
		updateFields["tel_id"] = poster.TelID
	}

	if poster.UserPhone != "" {
		updatedPosterModel.UserPhone = poster.UserPhone
		updateFields["user_phone"] = poster.UserPhone
	}

	if poster.Alert != "" { //todo modar fix alert
		if poster.Alert == "true" {
			updatedPosterModel.HasAlert = true
			updateFields["alert"] = true
		} else if poster.Alert == "false" {
			fmt.Println("modar alert is false")
			updatedPosterModel.HasAlert = false
			updateFields["alert"] = false
		}
	}

	if poster.Chat != "" { //todo modar fix chat
		if poster.Chat == "true" {
			updatedPosterModel.HasChat = true
			updateFields["chat"] = true
		} else if poster.Chat == "false" {
			updatedPosterModel.HasChat = false
			updateFields["chat"] = false
		}
	}

	if poster.Award != 0 {
		updatedPosterModel.Award = poster.Award
		updateFields["award"] = poster.Award
	}

	if poster.UserID != 0 {
		updatedPosterModel.UserID = poster.UserID
		updateFields["user_id"] = poster.UserID
	}

	if poster.State != "" {
		updatedPosterModel.State = poster.State
		updateFields["state"] = poster.State
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
		updateFields["images"] = poster.ImgUrls
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
	if len(addresses) != 0 {
		var esAddresses []DTO2.ESAddressDTO
		var updatedAddresses []Model2.Address
		for _, address := range addresses {
			esAddresses = append(esAddresses, DTO2.ESAddressDTO{
				Province:      address.Province,
				City:          address.City,
				AddressDetail: address.AddressDetail,
				Location: DTO2.Location{
					Latitude:  address.Latitude,
					Longitude: address.Longitude,
				},
			})
			updatedAddresses = append(updatedAddresses, Model2.Address{
				Province:      address.Province,
				City:          address.City,
				AddressDetail: address.AddressDetail,
				Latitude:      address.Latitude,
				Longitude:     address.Longitude,
			})
		}
		if strings.ToLower(role) == "user" {
			updatedPosterModel.State = "pending"
			updateFields["state"] = "pending"
		}

		r.db.Where("poster_id = ?", id).Delete(&Model2.Address{})
		_ = r.db.Model(&updatedPosterModel).Association("Addresses").Append(updatedAddresses)
		updateFields["addresses"] = esAddresses
	}

	//if len(poster.TagIds) != 0 {
	//	var updatedTags []Model2.Tag
	//	for _, tagId := range poster.TagIds {
	//		updatedTags = append(updatedTags, Model2.Tag{
	//			Url: url,
	//		})
	//	}
	//	_ = r.db.Model(&updatedPosterModel).Association("Images").Replace(updatedTags)
	//}

	result := r.db.Where("id = ?", id).Updates(updatedPosterModel)
	if result.RowsAffected == 0 {
		return errors.New("poster not found")
	}

	if result.Error != nil {
		return result.Error
	}

	updateFields["updated_at"] = time.Now()
	updateDoc := make(map[string]interface{})
	updateDoc["doc"] = updateFields

	err := esPosterCli.UpdatePoster(updateDoc, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PosterRepository) GetPosterByOwnerID(userID uint) (Model2.Poster, error) {
	var poster Model2.Poster
	r.db.First(&poster, "user_id = ?", userID)

	if poster.ID == 0 {
		return Model2.Poster{}, errors.New("poster not found")
	}

	return poster, nil
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

	var result = r.db.Preload("Poster").Limit(limit).Offset(offset).Order("created_at DESC")

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
	esPoster := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	var updatedPosterReportModel Model2.Poster
	updateFields := make(map[string]interface{})

	updatedPosterReportModel.State = state
	updateFields["state"] = state

	update := make(map[string]interface{})
	update["doc"] = updateFields

	result := r.db.Model(&Model2.Poster{}).Where("id = ?", id).Updates(updatedPosterReportModel)
	if result.Error != nil {
		return result.Error
	}

	err := esPoster.UpdatePoster(update, int(id))
	if err != nil {
		return err
	}

	return nil
}
