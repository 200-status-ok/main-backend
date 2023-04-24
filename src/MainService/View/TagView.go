package View

import (
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryView struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func CreateTagView(tag Model2.Category, c *gin.Context) {
	result := CategoryView{
		ID:   tag.ID,
		Name: tag.Name,
	}
	c.JSON(http.StatusOK, result)
}

func GetAllTagView(tags []Model2.Category, c *gin.Context) {
	var result []CategoryView
	for _, tag := range tags {
		result = append(result, CategoryView{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	c.JSON(http.StatusOK, result)
}
