package View

import (
	Model2 "github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type TagView struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func CreateTagView(tag Model2.Tag, c *gin.Context) {
	result := TagView{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt.UnixMilli(),
		UpdatedAt: tag.UpdatedAt.UnixMilli(),
	}
	c.JSON(http.StatusOK, result)
}

func GetAllTagView(tags []Model2.Tag, c *gin.Context) {
	var result []TagView
	for _, tag := range tags {
		result = append(result, TagView{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.UnixMilli(),
			UpdatedAt: tag.UpdatedAt.UnixMilli(),
		})
	}
	c.JSON(http.StatusOK, result)
}

type GeneratedPosterInfoView struct {
	Titles      []string `json:"titles"`
	Tags        []string `json:"tags"`
	Colors      []string `json:"colors"`
	Description string   `json:"description"`
}

func GeneratePosterInfoView(generatedTags dtos.GeneratedPosterTags, generatedColors dtos.GeneratedPosterColors, c *gin.Context) {
	var titlesResult []string
	for i := 0; i < int(math.Min(float64(len(generatedTags.Result.Tags)), 4)); i++ {
		titlesResult = append(titlesResult, generatedTags.Result.Tags[i].Tag.Fa)
	}

	var tagsResult []string
	for i := 0; i < int(math.Min(float64(len(generatedTags.Result.Tags)), 10)); i++ {
		tagsResult = append(tagsResult, generatedTags.Result.Tags[i].Tag.Fa)
	}

	var colorsResult []string
	for _, color := range generatedColors.Result.Colors.ForegroundColors {
		colorsResult = append(colorsResult, Utils.ColorParentsToPersian[color.ClosestPaletteColorParent])
	}

	description := " من یک " + tagsResult[0] + " با رنگ " + colorsResult[0] + " گم کرده ام "

	c.JSON(http.StatusOK, GeneratedPosterInfoView{
		Titles:      titlesResult,
		Tags:        tagsResult,
		Colors:      colorsResult,
		Description: description,
	})
}
