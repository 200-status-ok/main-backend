package main

import (
	"github.com/403-access-denied/main-backend/src/NotificationService/UseCase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	UseCase.SendToUser()
	r.Run(":8000")
}
