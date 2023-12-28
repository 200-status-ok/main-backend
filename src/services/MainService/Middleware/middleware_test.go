package Middleware

import (
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockTokenMaker struct {
	tokenMaker Token.Maker
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.True(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Equal(t, uint64(123), payload.(*Token.Payload).UserID, "Expected the correct user ID in the payload")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	jwtToken, _, err := token.MakeToken("09100570877", uint64(123), "user", time.Minute)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearer "+jwtToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.False(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Nil(t, payload, "Expected the payload to be nil")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.False(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Nil(t, payload, "Expected the payload to be nil")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearer "+"invalid.token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}

func TestAuthMiddleware_InvalidTypeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.True(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Equal(t, uint64(123), payload.(*Token.Payload).UserID, "Expected the correct user ID in the payload")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	jwtToken, _, err := token.MakeToken("09100570877", uint64(123), "user", time.Minute)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearerrr "+jwtToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}

func TestAdminAuthMiddleware_ValidAdminToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AdminAuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.True(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Equal(t, "Admin", payload.(*Token.Payload).Role, "Expected the correct role in the payload")
		assert.Equal(t, uint64(1234), payload.(*Token.Payload).UserID, "Expected the correct user ID in the payload")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	jwtToken, _, err := token.MakeToken("09941867752", uint64(1234), "Admin", time.Minute)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearer "+jwtToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

func TestAdminAuthMiddleware_ValidUserToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AdminAuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.True(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Equal(t, "Admin", payload.(*Token.Payload).Role, "Expected the correct role in the payload")
		assert.Equal(t, uint64(1234), payload.(*Token.Payload).UserID, "Expected the correct user ID in the payload")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	jwtToken, _, err := token.MakeToken("09941867752", uint64(1234), "user", time.Minute)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearer "+jwtToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}

func TestAdminAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AdminAuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.False(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Nil(t, payload, "Expected the payload to be nil")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}

func TestAdminAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AdminAuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.False(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Nil(t, payload, "Expected the payload to be nil")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearer "+"invalid.token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}

func TestAdminAuthMiddleware_InvalidTypeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	tokenMaker := mockTokenMaker{tokenMaker: token}
	router.Use(AdminAuthMiddleware(tokenMaker.tokenMaker))

	router.GET("/secured-endpoint", func(c *gin.Context) {
		payload, exists := c.Get(utils.AuthorizationPayloadKey)
		assert.True(t, exists, "Expected AuthorizationPayloadKey to exist in the context")
		assert.Equal(t, uint64(123), payload.(*Token.Payload).UserID, "Expected the correct user ID in the payload")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	jwtToken, _, err := token.MakeToken("09100570877", uint64(123), "Admin", time.Minute)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/secured-endpoint", nil)
	req.Header.Set(utils.AuthorizationHeaderKey, "Bearerrr "+jwtToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")
}
