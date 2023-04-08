package main

import (
	"github.com/403-access-denied/main-backend/docs"
	"github.com/403-access-denied/main-backend/src/Controller"
	"github.com/403-access-denied/main-backend/src/Token"
	"github.com/403-access-denied/main-backend/src/Utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Documentation for Golang web API(Gin framework)
// @version 1.0

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @BasePath /api/v1
func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	secretKey, _ := Utils.ReadFromEnvFile(".env", "JWT_SECRET")
	token, _ := Token.NewJWTMaker(secretKey)
	server := Controller.Server{Router: r, TokenMaker: token}
	server.MainController()

	r.Run(":8080")
}
