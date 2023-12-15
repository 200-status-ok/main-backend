package UseCase

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateSecretKeyForNewUser(t *testing.T) {
	userName := "09100570877"
	secretKey, err := GenerateSecretKeyForNewUser(userName)
	assert.NoError(t, err, "Expected no error")
	assert.NotEmpty(t, secretKey, "Expected a non-empty secret key")
}

func TestSendOTPResponse(t *testing.T) {
	t.Setenv("APP_ENV2", "test")
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/send-otp", SendOTPResponse)

	requestBody := `{ "username": "09100570877" }`
	w := performRequest(r, "POST", "/send-otp", requestBody)
	assert.Equal(t, 200, w.Code)
}

func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	if body != "" {
		req.Body = http.NoBody
		req.Body = ioutil.NopCloser(strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
