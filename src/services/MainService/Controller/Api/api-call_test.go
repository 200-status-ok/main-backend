package Api

import (
	"bytes"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"
)

func TestGeneratePosterInfo(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	router.GET("/api/v1/api-call/generate-poster-Info", GeneratePosterInfo)
	req1, err := http.NewRequest("GET",
		"/api/v1/api-call/generate-poster-Info?image_url=https://main-bucket.storage.iran.liara.space/20230524_154505_34ecbe928b029869.jpg",
		nil)
	assert.NoError(t, err)
	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	assert.Equal(t, http.StatusOK, W1.Code)
}

func TestImageUpload(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	router.POST("/api/v1/api-call/image-upload", ImageUpload)

	testFile, _ := os.Open("4.jpg")
	defer testFile.Close()
	fileHeader := mockMultipartFileHeader("4.jpg", len("4.jpg"), "image/jpeg")
	formFile := &multipart.FileHeader{
		Filename: "4.jpg",
		Header:   fileHeader.Header,
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", formFile.Filename)
	assert.NoError(t, err)
	io.Copy(part, testFile)
	writer.Close()

	req, err := http.NewRequest("POST", "/api/v1/api-call/image-upload", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	W := httptest.NewRecorder()
	router.ServeHTTP(W, req)

	fmt.Println(W.Body.String())
	assert.Equal(t, http.StatusOK, W.Code)
}

func TestUploadImageInServer(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	file, err := os.Open("4.jpeg")
	assert.NoError(t, err)
	defer file.Close()

	fmt.Println(file.Read([]byte{}))

	url, err := UseCase.UploadImageInServer(file, "4.jpeg", "liara", "main-bucket")
	assert.NoError(t, err)
	assert.NotEmpty(t, url)
	fmt.Println(url)
}

func mockMultipartFileHeader(filename string, size int, contentType string) *multipart.FileHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename))
	header.Set("Content-Type", contentType)
	fh := &multipart.FileHeader{
		Filename: filename,
		Header:   header,
		Size:     int64(size),
	}
	return fh
}
