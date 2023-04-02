package Controller

import (
	"github.com/403-access-denied/main-backend/src/Controller/Api"
	"github.com/403-access-denied/main-backend/src/Middleware"
	"github.com/403-access-denied/main-backend/src/Token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	TokenMaker Token.Maker
}

func (s *Server) MainController() {
	v1 := s.Router.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/send-otp", Api.SendOTP)
			user.POST("/login", Api.LoginUser)
			user.GET("/auth/google/login", Api.OAuth2Login)
			user.GET("/auth/google/callback", Api.GoogleCallback)
		}
		poster := v1.Group("/posters").Use(Middleware.AuthMiddleware(s.TokenMaker))
		{
			poster.GET("/", Api.GetPosters)
			poster.GET("/:id", Api.GetPoster)
			poster.POST("/", Api.CreatePoster)
			poster.PATCH("/:id", Api.UpdatePoster)
			poster.DELETE("/:id", Api.DeletePoster)
		}
	}

}
