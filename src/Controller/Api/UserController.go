package Api

import (
	"github.com/403-access-denied/main-backend/src/UseCase"
	"github.com/gin-gonic/gin"
)

// SendOTP LoginUser godoc
// @Summary send otp to user
// @Description send otp to user
// @Tags users
// @Accept  json
// @Produce  json
// @Param poster body UseCase.SendOTPRequest true "Send OTP"
// @Success 200 {object} UseCase.SendOTPRequest
// @Router /users/send-otp [post]
func SendOTP(c *gin.Context) {
	UseCase.SendOTPResponse(c)
}

// LoginUser godoc
// @Summary login user
// @Description login user
// @Tags users
// @Accept  json
// @Produce  json
// @Param poster body UseCase.VerifyOTPRequest true "Verify OTP"
// @Success 200 {object} View.UserView
// @Router /users/login [post]
func LoginUser(c *gin.Context) {
	UseCase.VerifyOtpResponse(c)
}
