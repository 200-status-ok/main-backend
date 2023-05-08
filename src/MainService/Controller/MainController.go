package Controller

import (
	Api2 "github.com/403-access-denied/main-backend/src/MainService/Controller/Api"
	"github.com/403-access-denied/main-backend/src/MainService/Middleware"
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	TokenMaker Token.Maker
	ChatWs     *Api2.ChatWS
}

func (s *Server) MainController() {
	v1 := s.Router.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			userAuthRoutes := user.Group("/authorize").Use(Middleware.AuthMiddleware(s.TokenMaker))
			{
				userAuthRoutes.GET("/", Api2.GetUser)
			}
			user.POST("/auth/otp/send", Api2.SendOTP)
			user.POST("/auth/otp/login", Api2.LoginUser)
			user.GET("/auth/google/login", Api2.OAuth2Login)
			user.GET("/auth/google/callback", Api2.GoogleCallback)
			user.GET("/", Api2.GetUsers)
			user.PATCH("/:id", Api2.UpdateUser)
			user.POST("/", Api2.CreateUser)
			user.DELETE("/:id", Api2.DeleteUser)
			user.GET("/payment/user_wallet", Api2.Payment)
			user.GET("/payment/user_wallet/:id", Api2.PaymentVerify)
		}
		// TODO add auth middleware
		poster := v1.Group("/posters")
		{
			poster.GET("/", Api2.GetPosters)
			poster.GET("/:id", Api2.GetPoster)
			poster.POST("/", Api2.CreatePoster)
			poster.PATCH("/:id", Api2.UpdatePoster)
			poster.DELETE("/:id", Api2.DeletePoster)
			poster.POST("/upload-image", Api2.UploadPosterImage)
		}
		report := v1.Group("/reports")
		{
			report.GET("/", Api2.GetPosterReports)
			report.GET("/:id", Api2.GetPosterReport)
			report.POST("/report-poster", Api2.CreatePosterReport)
			report.PATCH("/:id", Api2.UpdatePosterReport)
		}
		tags := v1.Group("/tags")
		{
			tags.GET("/", Api2.GetTags)
			tags.GET("/:id", Api2.GetTag)
			tags.PATCH("/:id", Api2.UpdateTag)
			tags.POST("/", Api2.CreateTag)
			tags.DELETE("/:id", Api2.DeleteTag)
		}
		api := v1.Group("/api-call")
		{
			api.GET("/generatePosterInfo", Api2.GeneratePosterInfo)
			api.POST("/predict", Api2.GetPhotoNSFWAi)
			api.GET("/predict-txt/", Api2.GetTextNSFW)
		}
		chats := v1.Group("/chats")
		{
			chats.GET("/join", s.ChatWs.JoinConversation)
			chats.POST("/conversation", Api2.CreateChatConversation)
		}
		admin := v1.Group("/admin")
		{
			admin.POST("/signup", Api2.SignupAdmin)
			admin.POST("/login", Api2.LoginAdmin)
		}
	}

}
