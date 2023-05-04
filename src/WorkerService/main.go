package main

import (
	"github.com/403-access-denied/main-backend/src/WorkerService/UseCase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	go UseCase.SendToUser()

	err := r.Run(":8000")
	if err != nil {
		return
	}
}
