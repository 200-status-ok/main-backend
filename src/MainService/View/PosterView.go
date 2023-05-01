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
	Categories  []Model2.Category   `json:"categories"`
	User        uint                `json:"user"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
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
			Categories:  poster.Categories,
			User:        poster.UserID,
			CreatedAt:   poster.CreatedAt,
			UpdatedAt:   poster.UpdatedAt,
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
		Categories:  poster.Categories,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
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
		Categories:  poster.Categories,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
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
		Categories:  poster.Categories,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
	}
	c.JSON(http.StatusOK, result)
}
