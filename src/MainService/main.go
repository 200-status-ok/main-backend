package main

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/Controller"
	"github.com/403-access-denied/main-backend/src/MainService/Controller/Api"
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/403-access-denied/main-backend/src/MainService/WebSocket"
	"github.com/403-access-denied/main-backend/src/MainService/docs"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// @title Swagger Documentation for Golang web API(Gin framework)
// @version 1.0
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @BasePath /api/v1
func main() {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Origin", "Allow-Origin", "Accept",
		"Accept-Encoding", "Accept-Language", "Connection", "Host", "Referer", "User-Agent",
		"Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-Site", "X-Requested-With", "X-Forwarded-For", "Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods", "Access-Control-Allow-Headers", "Access-Control-Allow-Credentials"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://71cb6a234a004aa492c0a6482b9a07e4@o4505154999025664.ingest.sentry.io/4505160636760064",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	r := gin.Default()
	r.Use(cors.New(config))
	r.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	secretKey := Utils.ReadFromEnvFile(".env", "JWT_SECRET")
	token, _ := Token.NewJWTMaker(secretKey)
	hub := WebSocket.NewHub()
	wsUseCase := Api.NewChatWS(hub)
	server := Controller.Server{Router: r, TokenMaker: token, ChatWs: wsUseCase}
	server.MainController()

	go wsUseCase.Hub.Run()

	r.Run(":8080")
}
