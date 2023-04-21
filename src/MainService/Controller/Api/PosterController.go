package Api

import (
	"github.com/403-access-denied/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// GetPosters godoc
// @Summary Get a list of all posters
// @Description Retrieves a list of all posters, sorted and paginated according to the given parameters
// @Tags posters
// @Accept  json
// @Produce  json
// @Param page_id query int true "Page ID" minimum(1) default(1)
// @Param page_size query int true "Page size" minimum(1) default(10)
// @Param sort query string false "Sort direction" enum(asc, desc) default(asc)
// @Param sort_by query string false "Sort by" enum(id, updated_at, created_at) default(created_at)
// @Param search_phrase query string false "Search phrase"
// @Param status query string false "Status" enum(lost, found, both) default(both)
// @Param time_start query int false "Time start"
// @Param time_end query int false "Time end"
// @Param only_rewards query bool false "Only rewards"
// @Param lat query float64 false "Latitude"
// @Param lon query float64 false "Longitude"
// @Param tag_ids query []int false "TagIds" collectionFormat(multi) example(1,2,3)
// @Success 200 {array} View.PosterView
// @Router /posters [get]
func GetPosters(c *gin.Context) {
	UseCase.GetPostersResponse(c)
}

// GetPoster godoc
// @Summary Get a poster by ID
// @Description Retrieves a poster by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200 {object} View.PosterView
// @Router /posters/{id} [get]
func GetPoster(c *gin.Context) {
	UseCase.GetPosterByIdResponse(c)
}

// CreatePoster godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags posters
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters [post]
func CreatePoster(c *gin.Context) {
	UseCase.CreatePosterResponse(c)
}

// UpdatePoster godoc
// @Summary Update a poster by ID
// @Description Updates a poster by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Param poster body UseCase.UpdatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters/{id} [patch]
func UpdatePoster(c *gin.Context) {
	UseCase.UpdatePosterResponse(c)
}

// DeletePoster godoc
// @Summary Delete a poster by ID
// @Description Deletes a poster by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200
// @Router /posters/{id} [delete]
func DeletePoster(c *gin.Context) {
	UseCase.DeletePosterByIdResponse(c)
}
