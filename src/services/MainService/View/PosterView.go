package View

import (
	Model2 "github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type PosterView struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Status      Model2.PosterStatus `json:"status"`
	Description string              `json:"description"`
	TelegramId  string              `json:"telegram_id"`
	UserPhone   string              `json:"user_phone"`
	Addresses   []AddressView       `json:"addresses"`
	Images      []ImageView         `json:"images"`
	Tags        []TagView           `json:"tags"`
	User        uint                `json:"user_id"`
	Award       float64             `json:"award"`
	CreatedAt   int64               `json:"created_at"`
	UpdatedAt   int64               `json:"updated_at"`
	State       string              `json:"state"`
	SpecialType string              `json:"special_type"`
}

type MarkedPosterView struct {
	ID        uint       `json:"id"`
	PosterID  uint       `json:"poster_id"`
	UserID    uint       `json:"user_id"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
	Poster    PosterView `json:"poster"`
}

type ImageView struct {
	ID        uint   `json:"id"`
	Url       string `json:"url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type AllPostersView struct {
	Total   int                 `json:"total"`
	MaxPage int                 `json:"max_page"`
	Posters []*dtos.ESPosterDTO `json:"posters"`
}

func GetPostersView(posters []*dtos.ESPosterDTO, totalPosters int, size int, c *gin.Context) {
	var result AllPostersView

	result.Total = totalPosters
	result.MaxPage = int(math.Ceil(float64(totalPosters) / float64(size)))
	result.Posters = posters

	c.JSON(http.StatusOK, result)
}

func GetPosterByIdView(poster Model2.Poster, c *gin.Context) {
	addressesView := make([]AddressView, 0)
	for _, address := range poster.Addresses {
		addressesView = append(addressesView, AddressView{
			ID:            address.ID,
			Province:      address.Province,
			City:          address.City,
			AddressDetail: address.AddressDetail,
			PosterID:      address.PosterId,
			Latitude:      address.Latitude,
			Longitude:     address.Longitude,
			CreatedAt:     address.CreatedAt.Unix(),
			UpdatedAt:     address.UpdatedAt.Unix(),
		})
	}
	imagesView := make([]ImageView, 0)
	for _, image := range poster.Images {
		imagesView = append(imagesView, ImageView{
			ID:        image.ID,
			Url:       image.Url,
			CreatedAt: image.CreatedAt.Unix(),
			UpdatedAt: image.UpdatedAt.Unix(),
		})
	}
	tagsView := make([]TagView, 0)
	for _, tag := range poster.Tags {
		tagsView = append(tagsView, TagView{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.Unix(),
			UpdatedAt: tag.UpdatedAt.Unix(),
		})
	}
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   addressesView,
		Images:      imagesView,
		TelegramId:  poster.TelegramID,
		UserPhone:   poster.UserPhone,
		Status:      poster.Status,
		Tags:        tagsView,
		User:        poster.UserID,
		Award:       poster.Award,
		CreatedAt:   poster.CreatedAt.Unix(),
		UpdatedAt:   poster.UpdatedAt.Unix(),
		State:       poster.State,
		SpecialType: poster.SpecialType,
	}
	c.JSON(http.StatusOK, result)
}

func CreatePosterView(poster Model2.Poster, c *gin.Context) {
	addressesView := make([]AddressView, 0)
	for _, address := range poster.Addresses {
		addressesView = append(addressesView, AddressView{
			ID:            address.ID,
			Province:      address.Province,
			City:          address.City,
			AddressDetail: address.AddressDetail,
			PosterID:      address.PosterId,
			Latitude:      address.Latitude,
			Longitude:     address.Longitude,
			CreatedAt:     address.CreatedAt.Unix(),
			UpdatedAt:     address.UpdatedAt.Unix(),
		})
	}
	imagesView := make([]ImageView, 0)
	for _, image := range poster.Images {
		imagesView = append(imagesView, ImageView{
			ID:        image.ID,
			Url:       image.Url,
			CreatedAt: image.CreatedAt.Unix(),
			UpdatedAt: image.UpdatedAt.Unix(),
		})
	}
	tagsView := make([]TagView, 0)
	for _, tag := range poster.Tags {
		tagsView = append(tagsView, TagView{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.Unix(),
			UpdatedAt: tag.UpdatedAt.Unix(),
		})
	}
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   addressesView,
		Images:      imagesView,
		Status:      poster.Status,
		Tags:        tagsView,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt.Unix(),
		UpdatedAt:   poster.UpdatedAt.Unix(),
		State:       poster.State,
		SpecialType: poster.SpecialType,
	}
	c.JSON(http.StatusOK, result)
}

type PosterReportView struct {
	ID          uint         `json:"id"`
	Poster      PosterView   `json:"poster"`
	Issuer      UserViewInfo `json:"issuer"`
	ReportType  string       `json:"report_type"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
}

func GetPosterReportsView(posterReports []Model2.PosterReport, c *gin.Context) {
	result := make([]PosterReportView, 0)
	for _, posterReport := range posterReports {
		addressesView := make([]AddressView, 0)
		for _, address := range posterReport.Poster.Addresses {
			addressesView = append(addressesView, AddressView{
				ID:            address.ID,
				Province:      address.Province,
				City:          address.City,
				AddressDetail: address.AddressDetail,
				PosterID:      address.PosterId,
				Latitude:      address.Latitude,
				Longitude:     address.Longitude,
				CreatedAt:     address.CreatedAt.Unix(),
				UpdatedAt:     address.UpdatedAt.Unix(),
			})
		}
		imagesView := make([]ImageView, 0)
		for _, image := range posterReport.Poster.Images {
			imagesView = append(imagesView, ImageView{
				ID:        image.ID,
				Url:       image.Url,
				CreatedAt: image.CreatedAt.Unix(),
				UpdatedAt: image.UpdatedAt.Unix(),
			})
		}
		tagsView := make([]TagView, 0)
		for _, tag := range posterReport.Poster.Tags {
			tagsView = append(tagsView, TagView{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt.Unix(),
				UpdatedAt: tag.UpdatedAt.Unix(),
			})
		}
		reportPoster := PosterView{
			ID:          posterReport.Poster.ID,
			Title:       posterReport.Poster.Title,
			Description: posterReport.Poster.Description,
			Addresses:   addressesView,
			Images:      imagesView,
			TelegramId:  posterReport.Poster.TelegramID,
			UserPhone:   posterReport.Poster.UserPhone,
			Status:      posterReport.Poster.Status,
			Tags:        tagsView,
			User:        posterReport.Poster.UserID,
			Award:       posterReport.Poster.Award,
			CreatedAt:   posterReport.Poster.CreatedAt.Unix(),
			UpdatedAt:   posterReport.Poster.UpdatedAt.Unix(),
			State:       posterReport.Poster.State,
			SpecialType: posterReport.Poster.SpecialType,
		}
		userInfo := UserViewInfo{
			Id:           posterReport.Issuer.ID,
			Username:     posterReport.Issuer.Username,
			Posters:      nil,
			MarkedPoster: nil,
			Wallet:       posterReport.Issuer.Wallet,
		}
		result = append(result, PosterReportView{
			ID:          posterReport.ID,
			Poster:      reportPoster,
			Issuer:      userInfo,
			ReportType:  posterReport.ReportType,
			Description: posterReport.Description,
			Status:      posterReport.Status,
		})
	}
	c.JSON(http.StatusOK, result)
}

func GetPosterReportByIdView(posterReport Model2.PosterReport, c *gin.Context) {
	addressesView := make([]AddressView, 0)
	for _, address := range posterReport.Poster.Addresses {
		addressesView = append(addressesView, AddressView{
			ID:            address.ID,
			Province:      address.Province,
			City:          address.City,
			AddressDetail: address.AddressDetail,
			PosterID:      address.PosterId,
			Latitude:      address.Latitude,
			Longitude:     address.Longitude,
			CreatedAt:     address.CreatedAt.Unix(),
			UpdatedAt:     address.UpdatedAt.Unix(),
		})
	}
	imagesView := make([]ImageView, 0)
	for _, image := range posterReport.Poster.Images {
		imagesView = append(imagesView, ImageView{
			ID:        image.ID,
			Url:       image.Url,
			CreatedAt: image.CreatedAt.Unix(),
			UpdatedAt: image.UpdatedAt.Unix(),
		})
	}
	tagsView := make([]TagView, 0)
	for _, tag := range posterReport.Poster.Tags {
		tagsView = append(tagsView, TagView{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.Unix(),
			UpdatedAt: tag.UpdatedAt.Unix(),
		})
	}
	reportPoster := PosterView{
		ID:          posterReport.Poster.ID,
		Title:       posterReport.Poster.Title,
		Description: posterReport.Poster.Description,
		Addresses:   addressesView,
		Images:      imagesView,
		TelegramId:  posterReport.Poster.TelegramID,
		UserPhone:   posterReport.Poster.UserPhone,
		Status:      posterReport.Poster.Status,
		Tags:        tagsView,
		User:        posterReport.Poster.UserID,
		Award:       posterReport.Poster.Award,
		CreatedAt:   posterReport.Poster.CreatedAt.Unix(),
		UpdatedAt:   posterReport.Poster.UpdatedAt.Unix(),
		State:       posterReport.Poster.State,
	}
	userInfo := UserViewInfo{
		Id:           posterReport.Issuer.ID,
		Username:     posterReport.Issuer.Username,
		Posters:      nil,
		MarkedPoster: nil,
		Wallet:       posterReport.Issuer.Wallet,
	}
	result := PosterReportView{
		ID:          posterReport.ID,
		Poster:      reportPoster,
		Issuer:      userInfo,
		ReportType:  posterReport.ReportType,
		Description: posterReport.Description,
		Status:      posterReport.Status,
	}
	c.JSON(http.StatusOK, result)
}
