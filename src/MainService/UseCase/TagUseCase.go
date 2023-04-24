package UseCase

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/View"
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
	tagRepository := Repository.NewCategoryRepository(DBConfiguration.GetDB())
	tags, err := tagRepository.CreateCategory(Model.Category{
		Name: tag.Name,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.CreateTagView(tags, c)
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=30"`
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
	tagRepository := Repository.NewCategoryRepository(DBConfiguration.GetDB())
	tags, err := tagRepository.UpdateCategory(id.ID, Model.Category{
		Name: tag.Name,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.CreateTagView(tags, c)
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
	tagRepository := Repository.NewCategoryRepository(DBConfiguration.GetDB())
	err := tagRepository.DeleteCategory(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	c.JSON(http.StatusOK, gin.H{})
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
	tagRepository := Repository.NewCategoryRepository(DBConfiguration.GetDB())
	tags, err := tagRepository.GetCategoryById(int(id.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.CreateTagView(tags, c)
}

func GetTagsResponse(c *gin.Context) {
	tagRepository := Repository.NewCategoryRepository(DBConfiguration.GetDB())
	tags, err := tagRepository.GetCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.GetAllTagView(tags, c)
}
