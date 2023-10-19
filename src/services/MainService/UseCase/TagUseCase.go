package UseCase

import (
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/View"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=30"`
}

func CreateTagResponse(c *gin.Context) {
	var tag CreateTagRequest
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tagRepository := Repository.NewCategoryRepository(pgsql.GetDB())
	tags, err := tagRepository.CreateCategory(Model.Tag{
		Name: tag.Name,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.CreateTagView(tags, c)
}

type UpdateTagRequest struct {
	Name  string `json:"name" binding:"max=30"`
	State string `json:"state" binding:"oneof=accepted rejected pending ''"`
}
type UpdateTagByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func UpdateTagByIdResponse(c *gin.Context) {
	var tag UpdateTagRequest
	var id UpdateTagByIdRequest
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagRepository := Repository.NewCategoryRepository(pgsql.GetDB())
	err := tagRepository.UpdateTag(id.ID, Model.Tag{
		Name:  tag.Name,
		State: tag.State,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag updated successfully"})
}

type DeleteTagByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func DeleteTagByIdResponse(c *gin.Context) {
	var id DeleteTagByIdRequest
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tagRepository := Repository.NewCategoryRepository(pgsql.GetDB())
	err := tagRepository.DeleteCategory(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//DBConfiguration.CloseDB()
	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

type GetTagByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func GetTagByIdResponse(c *gin.Context) {
	var id GetTagByIdRequest
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tagRepository := Repository.NewCategoryRepository(pgsql.GetDB())
	tags, err := tagRepository.GetCategoryById(int(id.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//DBConfiguration.CloseDB()
	View.CreateTagView(tags, c)
}

type GetTagsRequest struct {
	State string `form:"state,omitempty" binding:"omitempty,oneof=all pending accepted rejected"`
}

func GetTagsResponse(c *gin.Context) {
	var request GetTagsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagRepository := Repository.NewCategoryRepository(pgsql.GetDB())
	tags, err := tagRepository.GetTags(request.State)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetAllTagView(tags, c)
}

type GeneratePosterInfoRequest struct {
	ImageUrl string `form:"image_url"`
}

func GeneratePosterInfoResponse(c *gin.Context) {

	var request GeneratePosterInfoRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c1 := make(chan dtos.GeneratedPosterTags)
	c2 := make(chan dtos.GeneratedPosterColors)

	apiKey := "acc_52760e43313cc5e"
	apiSecret := "eab29250d009f3ae079998c5b3e6c83d"

	fmt.Println("Generating info...")
	fmt.Println("image url: " + request.ImageUrl)

	go func() {

		fmt.Println("modar Generating tags...")

		client := &http.Client{}

		reqUrl := "https://api.imagga.com/v2/tags?image_url=" +
			request.ImageUrl +
			"&language=fa" +
			"&limit=10" +
			"&threshold=15"

		req, _ := http.NewRequest("GET", reqUrl, nil)
		req.SetBasicAuth(apiKey, apiSecret)
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Error: Error when sending request to the server:", err)
			c1 <- dtos.GeneratedPosterTags{}
			return
		}

		defer resp.Body.Close()

		var generatedPosterTags dtos.GeneratedPosterTags
		err = json.NewDecoder(resp.Body).Decode(&generatedPosterTags)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			c1 <- dtos.GeneratedPosterTags{}
			return
		}

		fmt.Println("Generating Tags response:", generatedPosterTags)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c1 <- dtos.GeneratedPosterTags{}
			return
		}

		c1 <- generatedPosterTags
	}()

	go func() {

		fmt.Println("Generating Colors...")

		client := &http.Client{}

		reqUrl := "https://api.imagga.com/v2/colors?image_url=" +
			request.ImageUrl +
			"&extract_overall_colors=0"

		req, _ := http.NewRequest("GET", reqUrl, nil)
		req.SetBasicAuth(apiKey, apiSecret)
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Error: Error when sending request to the server:", err)
			c2 <- dtos.GeneratedPosterColors{}
			return
		}

		defer resp.Body.Close()

		var generatedPosterColors dtos.GeneratedPosterColors
		err = json.NewDecoder(resp.Body).Decode(&generatedPosterColors)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			c2 <- dtos.GeneratedPosterColors{}
			return
		}

		fmt.Println("Generating Colors response:", generatedPosterColors)

		c2 <- generatedPosterColors
	}()

	fmt.Println("Waiting for info...")
	generatedTagsResult := <-c1
	generatedColorsResult := <-c2

	fmt.Println("Info generated successfully!")

	View.GeneratePosterInfoView(generatedTagsResult, generatedColorsResult, c)
}
