package Utils

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
)

func UploadInArvanCloud(input multipart.File, fileName string) (string, error) {
	secretKey := utils.ReadFromEnvFile(".env", "ARVAN_SECRET_KEY")
	accessKey := utils.ReadFromEnvFile(".env", "ARVAN_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String("simin"),
		Endpoint:    aws.String("s3.ir-thr-at1.arvanstorage.ir"),
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("main-bucket"),
		Key:    aws.String(fileName),
		Body:   input,
	})
	if err != nil {
		return "", err
	}
	svc := s3.New(sess, &aws.Config{
		Region:   aws.String("simin"),
		Endpoint: aws.String("s3.ir-thr-at1.arvanstorage.ir"),
	})

	params := &s3.PutObjectAclInput{
		Bucket: aws.String("main-bucket"),
		Key:    aws.String(fileName),
		ACL:    aws.String("public-read"),
	}

	_, err = svc.PutObjectAcl(params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://main-bucket.s3.ir-thr-at1.arvanstorage.ir/%s", fileName), nil
}

func UploadInLiaraCloud(input multipart.File, fileName string) (string, error) {
	secretKey := utils.ReadFromEnvFile(".env", "LIARA_SECRET_KEY")
	accessKey := utils.ReadFromEnvFile(".env", "LIARA_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String("https://storage.iran.liara.space"),
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("main-bucket"),
		Key:    aws.String(fileName),
		Body:   input,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://main-bucket.storage.iran.liara.space/%s", fileName), nil
}
