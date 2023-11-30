package UseCase

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func ImageUploadResponse(c *gin.Context) {
	formHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}(file)

	uploadUrl, err := UploadImageInServer(file, newName, "arvan", "main-bucket")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": uploadUrl})
}

func UploadImageInServer(input multipart.File, fileName string, serverName string, bucketName string) (string, error) {
	var secretKey, accessKey, region, endPoint, resEndPoint string
	if serverName == "arvan" {
		secretKey = utils.ReadFromEnvFile(".env", "ARVAN_SECRET_KEY")
		accessKey = utils.ReadFromEnvFile(".env", "ARVAN_ACCESS_KEY")
		if os.Getenv("APP_ENV2") == "testing" {
			secretKey = "6bdaa1ac9227179b6a0d76c560b3563b97accc89"
			accessKey = "6c2c235b-ec0c-4f90-91cf-1c870c4cf4be"
		}
		region = "simin"
		endPoint = "s3.ir-thr-at1.arvanstorage.ir"
		resEndPoint = "https://main-bucket.s3.ir-thr-at1.arvanstorage.ir/"
	} else if serverName == "liara" {
		secretKey = utils.ReadFromEnvFile(".env", "LIARA_SECRET_KEY")
		accessKey = utils.ReadFromEnvFile(".env", "LIARA_ACCESS_KEY")
		if os.Getenv("APP_ENV2") == "testing" {
			secretKey = "ea642489-1d82-4a7f-97ee-0a63a6721748"
			accessKey = "fnenl8q3f2s9kk73"
		}
		region = "us-east-1"
		endPoint = "storage.iran.liara.space"
		resEndPoint = "https://main-bucket.storage.iran.liara.space/"
	} else {
		return "", fmt.Errorf("server name is not valid")
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      &region,
		Endpoint:    &endPoint,
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   input,
	})
	if serverName == "arvan" {
		svc := s3.New(sess, &aws.Config{
			Region:   aws.String(region),
			Endpoint: aws.String(endPoint),
		})

		params := &s3.PutObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
			ACL:    aws.String("public-read"),
		}
		_, err = svc.PutObjectAcl(params)
		if err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}
	return fmt.Sprintf(resEndPoint+"%s", fileName), nil
}
