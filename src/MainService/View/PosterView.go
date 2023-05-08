package View

import (
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PosterView struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Status      Model2.PosterStatus `json:"status"`
	Description string              `json:"description"`
	TelegramId  string              `json:"telegram_id"`
	UserPhone   string              `json:"phone_user"`
	Addresses   []Model2.Address    `json:"address"`
	Images      []Model2.Image      `json:"images"`
	Tags        []Model2.Tag        `json:"categories"`
	User        uint                `json:"user"`
	Award       float64             `json:"award"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	State       string              `json:"state"`
}

func GetPostersView(posters []Model2.Poster, c *gin.Context) {
	result := make([]PosterView, 0)
	for _, poster := range posters {
		result = append(result, PosterView{
			ID:          poster.ID,
			Title:       poster.Title,
			Description: poster.Description,
			Addresses:   poster.Addresses,
			Images:      poster.Images,
			Status:      poster.Status,
			Tags:        poster.Tags,
			User:        poster.UserID,
			Award:       poster.Award,
			CreatedAt:   poster.CreatedAt,
			UpdatedAt:   poster.UpdatedAt,
			State:       poster.State,
		})
	}
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
	}
	c.JSON(http.StatusOK, result)
}

func UpdatePosterView(poster Model2.Poster, c *gin.Context) {
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   poster.Addresses,
		Images:      poster.Images,
		Status:      poster.Status,
		TelegramId:  poster.TelegramID,
		UserPhone:   poster.UserPhone,
		Tags:        poster.Tags,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
		State:       poster.State,
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
