package Admin

import (
	"github.com/403-access-denied/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// GetUser godoc
// @Summary Get a User by ID
// @Description Retrieves a User by ID
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} View.UserViewID
// @Router /admin/user [get]
func GetUser(c *gin.Context) {
	UseCase.GetUserByIdResponse(c)
}

// GetUsers godoc
// @Summary Get a Users
// @Description Retrieves Users
// @Tags admin
// @Accept  json
// @Produce  json
// @Success 200 {array} View.UserViewID
// @Router /admin/users [get]
func GetUsers(c *gin.Context) {
	UseCase.GetUsersResponse(c)
}

// UpdateUser godoc
// @Summary Update a User by ID
// @Description Updates a User by ID
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body UseCase.UpdateUserRequest true "User"
// @Success 200 {object} View.UserViewIDs
// @Router /admin/user/{id} [patch]
func UpdateUser(c *gin.Context) {
	UseCase.UpdateUserByIdResponse(c)
}

// CreateUser godoc
// @Summary Create a User
// @Description Create a User
// @Tags admin
// @Accept  json
// @Produce  json
// @Param user body UseCase.CreateUserRequest true "User"
// @Success 200 {object} View.UserViewID
// @Router /admin/user [post]
func CreateUser(c *gin.Context) {
	UseCase.CreateUserResponse(c)
}

// DeleteUser godoc
// @Summary Delete a User by ID
// @Description Deletes a User by ID
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Router /admin/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	UseCase.DeleteUserByIdResponse(c)
}

// CreatePoster godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags admin
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /admin/poster [post]
func CreatePoster(c *gin.Context) {
	UseCase.CreatePosterResponse(c)
}

// UpdatePoster godoc
// @Summary Update a poster by ID
// @Description Updates a poster by ID
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Param poster body UseCase.UpdatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /admin/poster/{id} [patch]
func UpdatePoster(c *gin.Context) {
	UseCase.UpdatePosterResponse(c)
}

// DeletePoster godoc
// @Summary Delete a poster by ID
// @Description Deletes a poster by ID
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200
// @Router /admin/poster/{id} [delete]
func DeletePoster(c *gin.Context) {
	UseCase.DeletePosterByIdResponse(c)
}
