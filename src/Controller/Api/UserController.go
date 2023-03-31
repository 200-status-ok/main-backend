package Api

import (
	"github.com/403-access-denied/main-backend/src/UseCase"
	"github.com/gin-gonic/gin"
)

// LoginUser godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags posters
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters [post]

func Login(c *gin.Context) {
	UseCase.LoginResponse(c)
}

// VerifyOtp godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags posters
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters [post]

func VerifyOtp(c *gin.Context) {
	UseCase.VerifyOtpResponse(c)
}

// Validate godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags posters
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters [post]

func Validate(c *gin.Context) {
	UseCase.ValidateResponse(c)
}

// LogedIn godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags posters
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters [post]

func LogedIn(c *gin.Context) {
	UseCase.LogedInResponse(c)
}
