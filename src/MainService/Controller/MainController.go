package Controller

import (
	Api2 "github.com/403-access-denied/main-backend/src/MainService/Controller/Api"
	"github.com/403-access-denied/main-backend/src/MainService/Controller/Api/Admin"
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
		v1.Use(Middleware.SentryMiddleware())
		// todo replace the auth middleware with the admin middleware
		authorizeAdmin := v1.Group("/admin").Use(Middleware.AdminAuthMiddleware(s.TokenMaker))
		{
			authorizeAdmin.GET("/user/:userid", Admin.GetUser)
			authorizeAdmin.GET("/users", Admin.GetUsers)
			authorizeAdmin.PATCH("/user/:userid", Admin.UpdateUser)
			authorizeAdmin.POST("/user", Admin.CreateUser)
			authorizeAdmin.DELETE("/user/:userid", Admin.DeleteUser)

			authorizeAdmin.POST("/poster", Admin.CreatePoster)
			authorizeAdmin.PATCH("/poster/:id", Admin.UpdatePoster)
			authorizeAdmin.DELETE("/poster/:id", Admin.DeletePoster)
		}
		admin := v1.Group("/admin")
		{
			admin.POST("/signup", Admin.SignupAdmin)
			admin.POST("/login", Admin.LoginAdmin)
		}
		user := v1.Group("/users")
		{
			userAuthRoutes := user.Group("/authorize").Use(Middleware.AuthMiddleware(s.TokenMaker))
			{
				userAuthRoutes.GET("/", Api2.GetUser)
				userAuthRoutes.PATCH("/", Api2.UpdateUser)
				userAuthRoutes.DELETE("/", Api2.DeleteUser)
				userAuthRoutes.GET("/payment/user_wallet", Api2.Payment)
				userAuthRoutes.GET("/payment/user_wallet/verify", Api2.PaymentVerify)
				userAuthRoutes.GET("/payment/user_wallet/get_transactions", Api2.GetTransactions)
			}
			user.POST("/auth/otp/send", Api2.SendOTP)
			user.POST("/auth/otp/login", Api2.LoginUser)
			user.GET("/auth/google/login", Api2.OAuth2Login)
			user.GET("/auth/google/callback", Api2.GoogleCallback)
		}
		poster := v1.Group("/posters")
		{
			poster.GET("/", Api2.GetPosters)
			poster.GET("/:id", Api2.GetPoster)
			authPosters := poster.Group("/authorize").Use(Middleware.AuthMiddleware(s.TokenMaker))
			{
				authPosters.POST("/", Api2.CreatePoster)
				authPosters.PATCH("/:id", Api2.UpdatePoster)
				authPosters.DELETE("/:id", Api2.DeletePoster)
			}
			poster.POST("/upload-image", Api2.UploadPosterImage)
			poster.PATCH("/update-state", Api2.UpdatePosterState)
			poster.POST("/create-mock-data", Api2.CreateMockData)
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
			api.GET("/generate-poster-Info", Api2.GeneratePosterInfo)
			api.POST("/nsfw/photo", Api2.CheckPhotoNSFW)
			api.GET("/nsfw/text", Api2.CheckTextNSFW)
		}
		chats := v1.Group("/chats")
		{
			chats.GET("/join", s.ChatWs.JoinConversation)
			chatAuthorize := chats.Group("/authorize").Use(Middleware.AuthMiddleware(s.TokenMaker))
			{
				chatAuthorize.POST("/conversation", Api2.CreateConversation)
				chatAuthorize.GET("/conversations", Api2.AllUserConversations)
				chatAuthorize.GET("/conversation/:conversation_id", Api2.GetConversationById)
				chatAuthorize.GET("/history/:conversation_id", Api2.ConversationHistory)
			}
		}
	}

}
