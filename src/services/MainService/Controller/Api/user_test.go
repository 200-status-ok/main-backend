package Api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Middleware"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var OTP string
var JwtToken string
var TrackerID string

func TestSendOTP(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	router.POST("/api/v1/users/auth/otp/send", SendOTP)

	tests := []struct {
		name       string
		request    map[string]interface{}
		expected   map[string]interface{}
		expectCode int
	}{
		{
			name: "Good Case",
			request: map[string]interface{}{
				"username": "alifakhari@gmail.com",
			},
			expected: map[string]interface{}{
				"OTP": "",
			},
			expectCode: http.StatusOK,
		},
		{
			name:       "Bad Case - Missing Username",
			request:    map[string]interface{}{},
			expected:   nil,
			expectCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestJSON, _ := json.Marshal(test.request)
			req, err := http.NewRequest("POST", "/api/v1/users/auth/otp/send", bytes.NewBuffer(requestJSON))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			W := httptest.NewRecorder()
			router.ServeHTTP(W, req)

			if test.expected != nil {
				var response map[string]interface{}
				err = json.Unmarshal(W.Body.Bytes(), &response)
				assert.NoError(t, err)

				otpValue, ok := response["OTP"].(string)
				assert.True(t, ok, "Expected OTP to be a string")
				assert.NotEmpty(t, otpValue, "Expected OTP to be non-empty")
				OTP = otpValue
			}

			assert.Equal(t, test.expectCode, W.Code)
		})
	}
}

func TestVerifyOTP(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	router.POST("/api/v1/auth/otp/login", LoginUser)
	validRequest := map[string]interface{}{
		"username": "alifakhari@gmail.com",
		"otp":      OTP,
	}
	validRequestJSON, _ := json.Marshal(validRequest)
	req1, err := http.NewRequest("POST", "/api/v1/auth/otp/login",
		bytes.NewBuffer(validRequestJSON))
	assert.NoError(t, err)

	req1.Header.Set("Content-Type", "application/json")
	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	var response map[string]interface{}
	err = json.Unmarshal(W1.Body.Bytes(), &response)
	assert.NoError(t, err)

	tokenValue, ok := response["token"].(string)
	assert.True(t, ok, "Expected token to be a string")
	assert.NotEmpty(t, tokenValue, "Expected token to be non-empty")

	JwtToken = tokenValue

	fmt.Println(W1.Body.String())
	assert.Equal(t, http.StatusOK, W1.Code)
}

func TestGetUserByID(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/users/authorize", GetUser)
	req1, err := http.NewRequest("GET", "/api/v1/users/authorize", nil)
	req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
	assert.NoError(t, err)

	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	fmt.Println(W1.Body.String())
	assert.Equal(t, http.StatusOK, W1.Code)
}

func TestUpdateUserByID(t *testing.T) {
	t.Run("Good Case", func(t *testing.T) {
		t.Setenv("APP_ENV2", "testing")
		router := gin.Default()
		token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
		router.Use(Middleware.AuthMiddleware(token))
		router.PATCH("/api/v1/users/authorize", UpdateUser)
		randomEmail := Utils.EmailRandomGenerator()
		validRequest := map[string]interface{}{
			"username": randomEmail,
		}
		validRequestJSON, _ := json.Marshal(validRequest)
		req1, err := http.NewRequest("PATCH", "/api/v1/users/authorize", bytes.NewBuffer(validRequestJSON))
		req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
		assert.NoError(t, err)

		req1.Header.Set("Content-Type", "application/json")
		W1 := httptest.NewRecorder()
		router.ServeHTTP(W1, req1)

		fmt.Println(W1.Body.String())
		assert.Equal(t, http.StatusOK, W1.Code)
	})
	t.Run("Bad Case - Missing Authorization Token", func(t *testing.T) {
		t.Setenv("APP_ENV2", "testing")
		router := gin.Default()
		token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
		router.Use(Middleware.AuthMiddleware(token))
		router.PATCH("/api/v1/users/authorize", UpdateUser)
		validRequest := map[string]interface{}{
			"username": "",
		}
		validRequestJSON, _ := json.Marshal(validRequest)
		req1, err := http.NewRequest("PATCH", "/api/v1/users/authorize", bytes.NewBuffer(validRequestJSON))
		req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
		assert.NoError(t, err)

		req1.Header.Set("Content-Type", "application/json")
		W1 := httptest.NewRecorder()
		router.ServeHTTP(W1, req1)

		fmt.Println(W1.Body.String())
		assert.Equal(t, http.StatusBadRequest, W1.Code)
	})
}

func TestDeleteUserByID(t *testing.T) {
	t.Run("Good Case", func(t *testing.T) {
		t.Setenv("APP_ENV2", "testing")
		router := gin.Default()
		token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
		router.Use(Middleware.AuthMiddleware(token))
		router.DELETE("/api/v1/users/authorize/", DeleteUser)
		req1, err := http.NewRequest("DELETE", "/api/v1/users/authorize/", nil)
		req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
		assert.NoError(t, err)

		W1 := httptest.NewRecorder()
		router.ServeHTTP(W1, req1)

		fmt.Println(W1.Body.String())
		assert.Equal(t, http.StatusOK, W1.Code)
	})
	t.Run("Bad Case - Missing Authorization Token", func(t *testing.T) {
		t.Setenv("APP_ENV2", "testing")
		router := gin.Default()
		token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
		router.Use(Middleware.AuthMiddleware(token))
		router.DELETE("/api/v1/users/authorize/", DeleteUser)
		req1, err := http.NewRequest("DELETE", "/api/v1/users/authorize/", nil)
		assert.NoError(t, err)

		W1 := httptest.NewRecorder()
		router.ServeHTTP(W1, req1)

		fmt.Println(W1.Body.String())
		assert.Equal(t, http.StatusUnauthorized, W1.Code)
	})
}

func TestPayment(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/users/authorize/payment/user_wallet", Payment)
	req1, err := http.NewRequest("GET", "/api/v1/users/authorize/payment/user_wallet?url=http://localhost:8080/swagger/index.html&amount=1000", nil)
	req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
	assert.NoError(t, err)

	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	fmt.Println(W1.Body.String())
	var response map[string]interface{}
	err = json.Unmarshal(W1.Body.Bytes(), &response)
	assert.NoError(t, err)

	TrackerID = response["trackID"].(string)
	assert.Equal(t, http.StatusOK, W1.Code)

	time.Sleep(15 * time.Second)
}

func TestPaymentVerify(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/users/authorize/payment/user_wallet/verify", PaymentVerify)
	req1, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/authorize/payment/user_wallet/verify?track_id=%s", TrackerID), nil)
	req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
	assert.NoError(t, err)

	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	fmt.Println(W1.Body.String())
	assert.Equal(t, http.StatusOK, W1.Code)
}

func TestGetTransactions(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/users/authorize/payment/user_wallet/transactions", GetTransactions)
	req1, err := http.NewRequest("GET", "/api/v1/users/authorize/payment/user_wallet/transactions", nil)
	req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
	assert.NoError(t, err)

	W1 := httptest.NewRecorder()
	router.ServeHTTP(W1, req1)

	fmt.Println(W1.Body.String())
}
