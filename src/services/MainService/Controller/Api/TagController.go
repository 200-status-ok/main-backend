package Api

import (
	"github.com/200-status-ok/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// CreateTag godoc
// @Summary Create a Tag by ID
// @Description Creates a Tag by ID
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param tag body UseCase.CreateTagRequest true "Tag"
// @Success 200 {object} View.TagView
// @Router /tags [post]
func CreateTag(c *gin.Context) {
	UseCase.CreateTagResponse(c)
}

// UpdateTag godoc
// @Summary Update a Tag by ID
// @Description Updates a Tag by ID
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Param tag body UseCase.UpdateTagRequest true "Tag"
// @Success 200 {object} View.TagView
// @Router /tags/{id} [patch]
func UpdateTag(c *gin.Context) {
	UseCase.UpdateTagByIdResponse(c)
}

// DeleteTag godoc
// @Summary Delete a Tag by ID
// @Description Deletes a Tag by ID
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Success 200
// @Router /tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	UseCase.DeleteTagByIdResponse(c)
}

// GetTag godoc
// @Summary Get a Tag by ID
// @Description Retrieves a Tag by ID
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Success 200 {object} View.TagView
// @Router /tags/{id} [get]
func GetTag(c *gin.Context) {
	UseCase.GetTagByIdResponse(c)
}

// GetTags godoc
// @Summary Get all Tags
// @Description Retrieves Tags
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param state query string false "State" enum(all, accepted, rejected, pending) default(all)
// @Success 200 {array} View.TagView
// @Router /tags [get]
func GetTags(c *gin.Context) {
	UseCase.GetTagsResponse(c)
}
