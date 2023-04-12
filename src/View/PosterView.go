package View

import (
	"github.com/403-access-denied/main-backend/src/Model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PosterView struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Addresses   []Model.Address  `json:"address"`
	Images      []Model.Image    `json:"images"`
	Categories  []Model.Category `json:"categories"`
	User        uint             `json:"user"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func GetPostersView(posters []Model.Poster, c *gin.Context) {
	result := make([]PosterView, 0)
	for _, poster := range posters {
		result = append(result, PosterView{
			ID:          poster.ID,
			Title:       poster.Title,
			Description: poster.Description,
			Addresses:   poster.Addresses,
			Images:      poster.Images,
			Categories:  poster.Categories,
			User:        poster.UserID,
			CreatedAt:   poster.CreatedAt,
			UpdatedAt:   poster.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, result)
}

func GetPosterByIdView(poster Model.Poster, c *gin.Context) {
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   poster.Addresses,
		Images:      poster.Images,
		Categories:  poster.Categories,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
	}
	c.JSON(http.StatusOK, result)
}

func CreatePosterView(poster Model.Poster, c *gin.Context) {
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   poster.Addresses,
		Images:      poster.Images,
		Categories:  poster.Categories,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
	}
	c.JSON(http.StatusOK, result)
}

func UpdatePosterView(poster Model.Poster, c *gin.Context) {
	result := PosterView{
		ID:          poster.ID,
		Title:       poster.Title,
		Description: poster.Description,
		Addresses:   poster.Addresses,
		Images:      poster.Images,
		Categories:  poster.Categories,
		User:        poster.UserID,
		CreatedAt:   poster.CreatedAt,
		UpdatedAt:   poster.UpdatedAt,
	}
	c.JSON(http.StatusOK, result)
}
