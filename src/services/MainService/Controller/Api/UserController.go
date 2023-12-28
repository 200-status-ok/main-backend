package Api

import (
	"github.com/200-status-ok/main-backend/src/MainService/UseCase"
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

// GoogleLoginAndroid godoc
// @Summary login user with Google
// @Description login user with Google
// @Tags users
// @Accept  json
// @Produce  json
// @Param email query string true "Email"
// @Router /users/auth/google/login/android [get]
func GoogleLoginAndroid(c *gin.Context) {
	UseCase.GoogleLoginAndroidResponse(c)
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
// @Summary Get a User by Authorization Token
// @Description Retrieves a User by Authorization Token
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} View.UserViewInfo
// @Router /users/authorize [get]
func GetUser(c *gin.Context) {
	UseCase.GetUserByIdResponse(c)
}

// UpdateUser godoc
// @Summary Update a User by ID
// @Description Updates a User by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body UseCase.UpdateUserRequest true "User"
// @Success 200 {object} View.UserViewInfo
// @Router /users/authorize/ [patch]
func UpdateUser(c *gin.Context) {
	UseCase.UpdateUserByIdResponse(c)
}

// MarkPoster godoc
// @Summary Mark a Poster by ID
// @Description Marks a Poster by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200
// @Router /users/authorize/mark-poster/{poster_id} [patch]
func MarkPoster(c *gin.Context) {
	UseCase.MarkPosterResponse(c)
}

// UnmarkPoster godoc
// @Summary Unmark a Poster by ID
// @Description Unmarks a Poster by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200
// @Router /users/authorize/mark-poster/{poster_id} [delete]
func UnmarkPoster(c *gin.Context) {
	UseCase.UnmarkPosterResponse(c)
}

// DeleteUser godoc
// @Summary Delete a User by ID
// @Description Deletes a User by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200
// @Router /users/authorize/ [delete]
func DeleteUser(c *gin.Context) {
	UseCase.DeleteUserByIdResponse(c)
}

// Payment godoc
// @Summary Payment
// @Description Payment
// @Tags users
// @Accept  json
// @Produce  json
// @Param url query string true "URL"
// @Param amount query float64 true "Amount"
// @Success 200
// @Router /users/authorize/payment/user_wallet [get]
func Payment(c *gin.Context) {
	UseCase.PaymentResponse(c)
}

// PaymentVerify godoc
// @Summary Payment Verify
// @Description Payment Verify
// @Tags users
// @Accept  json
// @Produce  json
// @Param track_id query string true "Track ID"
// @Success 200
// @Router /users/authorize/payment/user_wallet/verify [get]
func PaymentVerify(c *gin.Context) {
	UseCase.PaymentVerifyResponse(c)
}

// GetTransactions godoc
// @Summary Get Transactions
// @Description Get Transactions
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} View.UserViewPayment
// @Router /users/authorize/payment/user_wallet/transactions [get]
func GetTransactions(c *gin.Context) {
	UseCase.GetTransactionsResponse(c)
}
