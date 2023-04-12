package UseCase

import (
	"github.com/403-access-denied/main-backend/src/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/DTO"
	"github.com/403-access-denied/main-backend/src/Repository"
	"github.com/403-access-denied/main-backend/src/View"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getPostersRequest struct {
	PageID   int    `form:"page_id" binding:"required,min=1,max=10"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=10"`
	Sort     string `form:"sort,omitempty" binding:"omitempty,oneof=asc desc"`
	SortBy   string `form:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at id"`
}

func GetPostersResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request getPostersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Sort = c.DefaultQuery("sort", "asc")
	request.SortBy = c.DefaultQuery("sort_by", "created_at")
	offset := (request.PageID - 1) * request.PageSize

	posters, err := posterRepository.GetAllPosters(request.PageSize, offset, request.Sort, request.SortBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.GetPostersView(posters, c)
}

type GetPosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func GetPosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request GetPosterByIdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poster, err := posterRepository.GetPosterById(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.GetPosterByIdView(poster, c)
}

type DeletePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func DeletePosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request DeletePosterByIdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := posterRepository.DeletePosterById(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	c.JSON(http.StatusOK, gin.H{"message": "Poster deleted"})
}

type CreatePosterRequest struct {
	Poster     DTO.PosterDTO
	Addresses  []DTO.AddressDTO
	ImgUrls    []string `json:"img_urls" binding:"required"`
	Categories []int    `json:"categories" binding:"required"`
}

func CreatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request CreatePosterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poster, err := posterRepository.CreatePoster(request.Poster, request.Addresses, request.ImgUrls, request.Categories)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.CreatePosterView(poster, c)
}

type UpdatePosterRequest struct {
	Poster     DTO.PosterDTO
	Addresses  []DTO.AddressDTO
	ImgUrls    []string `json:"img_urls" binding:"required"`
	Categories []int    `json:"categories" binding:"required"`
}

type UpdatePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func UpdatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request UpdatePosterRequest
	var id UpdatePosterByIdRequest
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poster, err := posterRepository.UpdatePoster(id.ID, request.Poster, request.Addresses, request.ImgUrls,
		request.Categories)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.UpdatePosterView(poster, c)
}
