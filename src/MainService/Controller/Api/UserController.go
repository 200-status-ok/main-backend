package Api

import (
	"github.com/403-access-denied/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// SendOTP LoginUser godoc
// @Summary send otp to user
// @Description send otp to user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body UseCase.SendOTPRequest true "Send OTP"
// @Success 200 {object} View.MessageView
// @Router /users/auth/otp/send [post]
func SendOTP(c *gin.Context) {
	UseCase.SendOTPResponse(c)
}

// LoginUser godoc
// @Summary login user
// @Description login user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body UseCase.VerifyOTPRequest true "Verify OTP"
// @Success 200 {object} View.UserView
// @Router /users/auth/otp/login [post]
func LoginUser(c *gin.Context) {
	UseCase.VerifyOtpResponse(c)
}

// OAuth2Login godoc
// @Summary login user with oauth2
// @Description login user with oauth2
// @Tags users
// @Accept  json
// @Produce  json
// @Router /users/auth/google/login [get]
func OAuth2Login(c *gin.Context) {
	UseCase.OAuth2LoginResponse(c)
}

// GoogleCallback godoc
// @Summary google callback
// @Description google callback
// @Tags users
// @Accept  json
// @Produce  json
// @Router /users/auth/google/callback [get]
func GoogleCallback(c *gin.Context) {
	UseCase.GoogleCallbackResponse(c)
}

// GetUser godoc
// @Summary Get a User by ID
// @Description Retrieves a User by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} View.UserViewID
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	UseCase.GetUserByIdResponse(c)
}

// GetUsers godoc
// @Summary Get a Users
// @Description Retrieves Users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} View.UserViewID
// @Router /users [get]
func GetUsers(c *gin.Context) {
	UseCase.GetUsersResponse(c)
}

// UpdateUser godoc
// @Summary Update a User by ID
// @Description Updates a User by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body UseCase.UpdateUserRequest true "User"
// @Success 200 {object} View.UserViewIDs
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	UseCase.UpdateUserByIdResponse(c)
}

// CreateUser godoc
// @Summary Create a User
// @Description Create a User
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body UseCase.CreateUserRequest true "User"
// @Success 200 {object} View.UserViewID
// @Router /users [post]
func CreateUser(c *gin.Context) {
	UseCase.CreateUserResponse(c)
}

// DeleteUser godoc
// @Summary Delete a User by ID
// @Description Deletes a User by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	UseCase.DeleteUserByIdResponse(c)
}
