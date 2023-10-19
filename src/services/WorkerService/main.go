package main

import (
	"github.com/200-status-ok/main-backend/src/WorkerService/UseCase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	go UseCase.SendToUser()
	go UseCase.CheckPhotoNSFW()
	go UseCase.CheckTagNSFW()

	err := r.Run(":8081")
	if err != nil {
		return
	}
}
