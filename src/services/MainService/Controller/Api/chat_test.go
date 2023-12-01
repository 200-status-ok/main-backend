package Api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Middleware"
	"github.com/200-status-ok/main-backend/src/MainService/RealtimeChat"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type chatServer struct {
	chatWS *ChatWS
}

var hub = RealtimeChat.NewHub()
var wsUseCase = NewChatWS(hub)
var chat = chatServer{
	chatWS: wsUseCase,
}
var jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjoiMjAyMy0xMi0wN1QwODo0OTo0MC43NzQ3NzExNzZaIiwiaWQiOiJkNjA5MmE2My0zYzZjLTQ2YWYtYmY3MS1mMjI3NjMxYjBmNjQiLCJpc3N1ZWRBdCI6IjIwMjMtMTEtMzBUMDg6NDk6NDAuNzc0NzcwMzU4WiIsInJvbGUiOiJVc2VyIiwidXNlcklkIjo1MCwidXNlcm5hbWUiOiJhbGlmYWtoYXJpQGdtYWlsLmNvbSJ9.PKG9Sn6iC3sXya5Q-9_uoCMb20NhMf4KVkBKbvdBDkg"

func TestChatWS_SendMessage(t *testing.T) {
	t.Run("Exist Conversation", func(t *testing.T) {
		t.Setenv("APP_ENV2", "testing")
		router := gin.Default()
		token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
		router.Use(Middleware.AuthMiddleware(token))
		router.POST("/api/v1/chat/authorize/message", chat.chatWS.SendMessage)
		validReq := map[string]interface{}{
			"content":         "Hello",
			"conversation_id": 3,
			"post_id":         1,
			"receiver_id":     37,
			"sender_id":       38,
			"type":            "text",
		}
		requestJSON, _ := json.Marshal(validReq)
		req, err := http.NewRequest("POST", "/api/v1/chat/authorize/message", bytes.NewBuffer(requestJSON))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		W := httptest.NewRecorder()
		router.ServeHTTP(W, req)
		assert.Equal(t, http.StatusOK, W.Code)
		fmt.Println(W.Body.String())
	})
	t.Run("Not Exist Conversation", func(t *testing.T) {
		t.Setenv("APP_ENV2", "testing")
		router := gin.Default()
		token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
		router.Use(Middleware.AuthMiddleware(token))
		router.POST("/api/v1/chat/authorize/message", chat.chatWS.SendMessage)
		validReq := map[string]interface{}{
			"content":         "Hello, how are you?",
			"conversation_id": -1,
			"post_id":         1,
			"receiver_id":     37,
			"sender_id":       50,
			"type":            "text",
		}
		requestJSON, _ := json.Marshal(validReq)
		req, err := http.NewRequest("POST", "/api/v1/chat/authorize/message", bytes.NewBuffer(requestJSON))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		W := httptest.NewRecorder()
		router.ServeHTTP(W, req)
		assert.Equal(t, http.StatusBadRequest, W.Code)
		fmt.Println(W.Body.String())
	})
}

func TestChatWS_ReadMessages(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.POST("/api/v1/chat/authorize/read", chat.chatWS.ReadMessages)
	validReq := map[string]interface{}{
		"sender_id": 38,
		"message_ids": []int{
			22,
			23,
		},
	}
	requestJSON, _ := json.Marshal(validReq)
	req, err := http.NewRequest("POST", "/api/v1/chat/authorize/read", bytes.NewBuffer(requestJSON))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	W := httptest.NewRecorder()
	router.ServeHTTP(W, req)
	assert.Equal(t, http.StatusOK, W.Code)
	fmt.Println(W.Body.String())
}

func TestAllUserConversations(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/chat/authorize/conversation", AllUserConversations)

	req, err := http.NewRequest("GET", "/api/v1/chat/authorize/conversation", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	W := httptest.NewRecorder()
	router.ServeHTTP(W, req)
	assert.Equal(t, http.StatusOK, W.Code)
}

func TestGetConversationById(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/chat/authorize/conversation/:conversation_id", GetConversationById)

	req, err := http.NewRequest("GET", "/api/v1/chat/authorize/conversation/3", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	W := httptest.NewRecorder()
	router.ServeHTTP(W, req)
	assert.Equal(t, http.StatusOK, W.Code)
}

func TestConversationHistory(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.GET("/api/v1/chat/authorize/history/:conversation_id", ConversationHistory)

	req, err := http.NewRequest("GET", "/api/v1/chat/authorize/history/3?page_id=1&page_size=10", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	W := httptest.NewRecorder()
	router.ServeHTTP(W, req)
	assert.Equal(t, http.StatusOK, W.Code)
	fmt.Println(W.Body.String())
}

func TestUpdateConversation(t *testing.T) {
	t.Setenv("APP_ENV2", "testing")
	router := gin.Default()
	token, _ := Token.NewJWTMaker("qwertyuiopasdfghjklzxcvbnm123456")
	router.Use(Middleware.AuthMiddleware(token))
	router.PATCH("/api/v1/chat/authorize/conversation/:conversation_id", UpdateConversation)

	validReq := map[string]interface{}{
		"name": "conversation3",
	}
	requestJSON, _ := json.Marshal(validReq)
	req, err := http.NewRequest("PATCH", "/api/v1/chat/authorize/conversation/3", bytes.NewBuffer(requestJSON))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	W := httptest.NewRecorder()
	router.ServeHTTP(W, req)
	assert.Equal(t, http.StatusOK, W.Code)
}
