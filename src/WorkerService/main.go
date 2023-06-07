package main

import (
	"github.com/403-access-denied/main-backend/src/WorkerService/UseCase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	go UseCase.SendToUser()
	go UseCase.CheckPhotoNSFW()
	go UseCase.CheckTagNSFW()

	err := r.Run(":8000")
	if err != nil {
		return
	}
}
