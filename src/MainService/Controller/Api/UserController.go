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
