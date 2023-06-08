package View

import (
	"github.com/403-access-denied/main-backend/src/MainService/DTO"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

type PosterView struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Status      Model2.PosterStatus `json:"status"`
	Description string              `json:"description"`
	TelegramId  string              `json:"telegram_id"`
	UserPhone   string              `json:"user_phone"`
	Addresses   []Model2.Address    `json:"addresses"`
	Images      []Model2.Image      `json:"images"`
	Tags        []Model2.Tag        `json:"tags"`
	User        uint                `json:"user_id"`
	Award       float64             `json:"award"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	State       string              `json:"state"`
	SpecialType string              `json:"special_type"`
}

type AllPostersView struct {
	Total   int                `json:"total"`
	MaxPage int                `json:"max_page"`
	Posters []*DTO.ESPosterDTO `json:"posters"`
}

func GetPostersView(posters []*DTO.ESPosterDTO, totalPosters int, size int, c *gin.Context) {
	var result AllPostersView

	result.Total = totalPosters
	result.MaxPage = int(math.Ceil(float64(totalPosters) / float64(size)))
	result.Posters = posters

	c.JSON(http.StatusOK, result)
}

func GetPosterByIdView(poster Model2.Poster, c *gin.Context) {
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   poster.Addresses,
		Images:      poster.Images,
		TelegramId:  poster.TelegramID,
		UserPhone:   poster.UserPhone,
		Status:      poster.Status,
		Tags:        poster.Tags,
		User:        poster.UserID,
		Award:       poster.Award,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
		State:       poster.State,
		SpecialType: poster.SpecialType,
	}
	c.JSON(http.StatusOK, result)
}

func CreatePosterView(poster Model2.Poster, c *gin.Context) {
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   poster.Addresses,
		Images:      poster.Images,
		Status:      poster.Status,
		Tags:        poster.Tags,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
		State:       poster.State,
		SpecialType: poster.SpecialType,
	}
	c.JSON(http.StatusOK, result)
}

type PosterReportView struct {
	ID          uint          `json:"id"`
	Poster      Model2.Poster `json:"poster"`
	Issuer      Model2.User   `json:"issuer"`
	ReportType  string        `json:"report_type"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
}

func GetPosterReportsView(posterReports []Model2.PosterReport, c *gin.Context) {
	result := make([]PosterReportView, 0)
	for _, posterReport := range posterReports {
		result = append(result, PosterReportView{
			ID:          posterReport.ID,
			Poster:      posterReport.Poster,
			Issuer:      posterReport.Issuer,
			ReportType:  posterReport.ReportType,
			Description: posterReport.Description,
			Status:      posterReport.Status,
		})
	}
	c.JSON(http.StatusOK, result)
}

func GetPosterReportByIdView(posterReport Model2.PosterReport, c *gin.Context) {
	result := PosterReportView{
		ID:          posterReport.ID,
		Poster:      posterReport.Poster,
		Issuer:      posterReport.Issuer,
		ReportType:  posterReport.ReportType,
		Description: posterReport.Description,
		Status:      posterReport.Status,
	}
	c.JSON(http.StatusOK, result)
}
