package main

import (
	"github.com/403-access-denied/main-backend/docs"
	"github.com/403-access-denied/main-backend/src/Controller"
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
	// set /api/v1 as base path for all routes
	r.Group("/api/v1")
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Controller.MainController(r)
	r.Run(":8080")
}
