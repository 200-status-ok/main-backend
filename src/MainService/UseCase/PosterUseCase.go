package UseCase

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	DTO2 "github.com/403-access-denied/main-backend/src/MainService/DTO"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/View"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getPostersRequest struct {
	PageID       int     `form:"page_id" binding:"required,min=1"`
	PageSize     int     `form:"page_size" binding:"required,min=5,max=10"`
	Sort         string  `form:"sort,omitempty" binding:"omitempty,oneof=asc desc"`
	SortBy       string  `form:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at id"`
	Status       string  `form:"status,omitempty" binding:"oneof=lost found both"`
	SearchPhrase string  `form:"search_phrase,omitempty"`
	TimeStart    int64   `form:"time_start,omitempty"`
	TimeEnd      int64   `form:"time_end,omitempty"`
	onlyRewards  bool    `form:"only_rewards,omitempty" binding:"oneof=true false"`
	Lat          float64 `form:"lat,omitempty"`
	Lon          float64 `form:"lon,omitempty"`
}

func GetPostersResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())

	var request getPostersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset := (request.PageID - 1) * request.PageSize
	request.Sort = c.DefaultQuery("sort", "asc")
	request.SortBy = c.DefaultQuery("sort_by", "created_at")
	request.Status = c.DefaultQuery("status", "both")
	request.SearchPhrase = c.DefaultQuery("search_phrase", "")
	//todo add other fields

	filterObject := DTO2.FilterObject{
		Status:       request.Status,
		SearchPhrase: request.SearchPhrase,
		TimeStart:    request.TimeStart,
		TimeEnd:      request.TimeEnd,
		OnlyRewards:  request.onlyRewards,
		Lat:          request.Lat,
		Lon:          request.Lon,
	}

	posters, err := posterRepository.GetAllPosters(request.PageSize, offset, request.Sort, request.SortBy, filterObject)

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
	Poster     DTO2.PosterDTO
	Addresses  []DTO2.AddressDTO
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
	Poster     DTO2.PosterDTO
	Addresses  []DTO2.AddressDTO
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
