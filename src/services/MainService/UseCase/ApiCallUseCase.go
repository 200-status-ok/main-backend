package UseCase

import (
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"
)

func ImageUploadResponse(c *gin.Context) {
	formHeader, err := c.FormFile("image")
	fileName := formHeader.Filename
	extension := path.Ext(fileName)

	currentTime := time.Now().Format("20060102_150405")
	randomString := strconv.FormatInt(rand.Int63(), 16)
	newName := currentTime + "_" + randomString + extension
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := formHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	uploadUrl, err := Utils.UploadInLiaraCloud(file, newName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": uploadUrl})
}
