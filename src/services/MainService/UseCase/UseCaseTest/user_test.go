package UseCaseTest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Controller/Api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendOTP(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	router.POST("/api/v1/auth/otp/send", Api.SendOTP)
	validRequest := map[string]interface{}{
		"username": "alifakhary622@gmail.com",
	}
	validRequestJSON, _ := json.Marshal(validRequest)
	req1, err := http.NewRequest("POST", "/api/v1/auth/otp/send",
		bytes.NewBuffer(validRequestJSON))
	assert.NoError(t, err)

	req1.Header.Set("Content-Type", "application/json")

	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	fmt.Println(W1.Body.String())
	assert.Equal(t, http.StatusOK, W1.Code)
}

func TestVerifyOTP(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	router.POST("/api/v1/auth/otp/login", Api.LoginUser)
	validRequest := map[string]interface{}{
		"username": "alifakhary622@gmail.com",
		"otp":      "246981",
	}
	validRequestJSON, _ := json.Marshal(validRequest)
	req1, err := http.NewRequest("POST", "/api/v1/auth/otp/login",
		bytes.NewBuffer(validRequestJSON))
	assert.NoError(t, err)

	req1.Header.Set("Content-Type", "application/json")
	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	fmt.Println(W1.Body.String())
	assert.Equal(t, http.StatusOK, W1.Code)
}
