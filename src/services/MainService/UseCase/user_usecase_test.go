package UseCase

import (
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Middleware"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
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

var OTP string
var JwtToken string

func TestSendOTPResponse(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/send-otp", SendOTPResponse)

	tests := []struct {
		name       string
		request    string
		expectCode int
		appEnv     string
	}{
		{
			name:       "Good Case",
			request:    `{ "username": "09100570877" }`,
			expectCode: http.StatusOK,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Invalid Username",
			request:    `{ "username": "0dekfoeo" }`,
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Username",
			request:    `{}`,
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Username",
			request:    `{}`,
			expectCode: http.StatusBadRequest,
			appEnv:     "production",
		},
		{
			name:       "Bad Case - Missing Username",
			request:    `{}`,
			expectCode: http.StatusBadRequest,
			appEnv:     "development",
		},
		{
			name:       "Good Case",
			request:    `{ "username": "09100570877" }`,
			expectCode: http.StatusBadRequest,
			appEnv:     "development",
		},
		{
			name:       "Good Case",
			request:    `{ "username": "09100570877" }`,
			expectCode: http.StatusBadRequest,
			appEnv:     "production",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("APP_ENV2", test.appEnv)
			w := performRequest(r, "POST", "/send-otp", test.request)
			assert.Equal(t, test.expectCode, w.Code)
			if test.expectCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				otpValue, ok := response["OTP"].(string)
				assert.True(t, ok, "Expected OTP to be a string")
				assert.NotEmpty(t, otpValue, "Expected OTP to be non-empty")
				OTP = otpValue
			}
		})
	}
}

func TestVerifyOTPResponse(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/verify-otp", VerifyOtpResponse)

	tests := []struct {
		name       string
		request    string
		expectCode int
		appEnv     string
	}{
		{
			name:       "Good Case",
			request:    fmt.Sprintf(`{ "username": "09100570877", "otp": "%s" }`, OTP),
			expectCode: http.StatusOK,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Invalid Username",
			request:    fmt.Sprintf(`username": "fjfoekpl", "otp": "%s" }`, OTP),
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Username",
			request:    fmt.Sprintf(`{ "otp": "%s" }`, OTP),
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing OTP",
			request:    fmt.Sprintf(`{ "username": "09100570877" }`),
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Username and OTP",
			request:    `{}`,
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Username and OTP",
			request:    `{}`,
			expectCode: http.StatusBadRequest,
			appEnv:     "production",
		},
		{
			name:       "Bad Case - Missing Username and OTP",
			request:    `{}`,
			expectCode: http.StatusBadRequest,
			appEnv:     "development",
		},
		{
			name:       "Good Case",
			request:    fmt.Sprintf(`{ "username": "09100570877", "otp": "%s" }`, OTP),
			expectCode: http.StatusBadRequest,
			appEnv:     "development",
		},
		{
			name:       "Good Case",
			request:    fmt.Sprintf(`{ "username": "09100570877", "otp": "%s" }`, OTP),
			expectCode: http.StatusBadRequest,
			appEnv:     "production",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("APP_ENV2", test.appEnv)
			w := performRequest(r, "POST", "/verify-otp", test.request)
			assert.Equal(t, test.expectCode, w.Code)

			if test.expectCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				tokenValue, ok := response["token"].(string)
				assert.True(t, ok, "Expected token to be a string")
				assert.NotEmpty(t, tokenValue, "Expected token to be non-empty")

				JwtToken = tokenValue
			}
		})
	}
}

func TestGoogleLoginAndroidResponse(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/google-login-android", GoogleLoginAndroidResponse)

	// get request that sent logged in email in query
	tests := []struct {
		name       string
		request    string
		expectCode int
		appEnv     string
	}{
		{
			name:       "Good Case",
			request:    `?email=alifakhary622@gmail.com`,
			expectCode: http.StatusOK,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Invalid Email",
			request:    `? email=alifakhary622gmail.com`,
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Email",
			request:    ``,
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing Email",
			request:    ``,
			expectCode: http.StatusBadRequest,
			appEnv:     "production",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("APP_ENV2", test.appEnv)
			w := performRequest(r, "GET", "/google-login-android"+test.request, "")
			assert.Equal(t, test.expectCode, w.Code)
		})
	}
}

func TestGetUserByIdResponse(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/users/authorize", GetUserByIdResponse)

	// get users by JWT token in header
	tests := []struct {
		name       string
		request    string
		expectCode int
		appEnv     string
	}{
		{
			name:       "Good Case",
			request:    "",
			expectCode: http.StatusOK,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Invalid JWT Token",
			request:    "",
			expectCode: http.StatusUnauthorized,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Missing JWT Token",
			request:    "",
			expectCode: http.StatusUnauthorized,
			appEnv:     "testing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("APP_ENV2", test.appEnv)
			req1, err := http.NewRequest("GET", "/api/v1/users/authorize", nil)
			if test.name == "Good Case" {
				req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))
			} else if test.name == "Bad Case - Invalid JWT Token" {
				req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalid token"))
			} else if test.name == "Bad Case - Missing JWT Token" {
				req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ""))
			}
			assert.NoError(t, err)

			W1 := httptest.NewRecorder()
			router.ServeHTTP(W1, req1)

			fmt.Println(W1.Body.String())
			assert.Equal(t, test.expectCode, W1.Code)
		})
	}
}

func TestMarkPosterResponse(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	r.Use(Middleware.AuthMiddleware(token))
	r.PATCH("/mark-poster/:poster_id", MarkPosterResponse)

	// get request that sent poster id in path
	tests := []struct {
		name       string
		request    string
		expectCode int
		appEnv     string
	}{
		{
			name:       "Good Case",
			request:    `1`,
			expectCode: http.StatusOK,
			appEnv:     "testing",
		},
		{
			name:       "Bad Case - Invalid Poster ID",
			request:    `0`,
			expectCode: http.StatusBadRequest,
			appEnv:     "testing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("APP_ENV2", test.appEnv)
			req1, err := http.NewRequest("PATCH", "/mark-poster/"+test.request, nil)
			assert.NoError(t, err)
			req1.Header.Set("Authorization", fmt.Sprintf("Bearer %s", JwtToken))

			W1 := httptest.NewRecorder()
			r.ServeHTTP(W1, req1)
			assert.Equal(t, test.expectCode, W1.Code)
		})
	}
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
