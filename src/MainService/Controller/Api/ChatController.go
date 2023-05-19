package Api

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/403-access-denied/main-backend/src/MainService/View"
	"github.com/403-access-denied/main-backend/src/MainService/WebSocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatWS struct {
	Hub *WebSocket.Hub
}

func NewChatWS(hub *WebSocket.Hub) *ChatWS {
	return &ChatWS{Hub: hub}
}

type JoinConversationReq struct {
	ConversationID uint   `form:"conv_id" binding:"required"`
	Token          string `form:"token" binding:"required"`
}

// JoinConversation JoinChat godoc
// @Summary JoinConversation a chat room
// @Description JoinConversation a chat room
// @Tags Chat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Message body WebSocket.MessageWithType true "Message"
// @Param conv_id query uint true "Conversation ID"
// @Param token query string true "Token"
// @Success 200 {object} string
// @Router /chats/join [get]
func (wsUseCase *ChatWS) JoinConversation(c *gin.Context) {
	var request JoinConversationReq
	secretKey := Utils.ReadFromEnvFile(".env", "JWT_SECRET")
	tokenMaker, _ := Token.NewJWTMaker(secretKey)
	//payload := c.MustGet("authorization_payload").(*Token.Payload)
	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := request.Token
	payload, err := tokenMaker.VerifyToken(token)
	conversation, err := chatRepo.GetConversationById(request.ConversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if _, ok := wsUseCase.Hub.ChatConversation[int(conversation.ID)]; !ok {
		wsUseCase.Hub.ChatConversation[int(conversation.ID)] = &WebSocket.ConversationChat{
			ID:     int(conversation.ID),
			Name:   fmt.Sprintf("conversation-%d", conversation.ID),
			Owner:  &WebSocket.Client{},
			Member: &WebSocket.Client{},
		}
		memberClient := &WebSocket.Client{
			Conn:           &websocket.Conn{},
			Message:        make(chan *WebSocket.Message, 100),
			ID:             int(conversation.MemberID),
			Role:           WebSocket.Member,
			ConversationID: int(conversation.ID),
			IsConnected:    false,
		}
		wsUseCase.Hub.ChatConversation[int(conversation.ID)].Member = memberClient
		ownerClient := &WebSocket.Client{
			Conn:           &websocket.Conn{},
			Message:        make(chan *WebSocket.Message, 15),
			ID:             int(conversation.OwnerID),
			Role:           WebSocket.Owner,
			ConversationID: int(conversation.ID),
			IsConnected:    false,
		}
		wsUseCase.Hub.ChatConversation[int(conversation.ID)].Owner = ownerClient
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var client *WebSocket.Client
	if uint(payload.UserID) == conversation.OwnerID {
		client = wsUseCase.Hub.ChatConversation[int(conversation.ID)].Owner
	} else {
		client = wsUseCase.Hub.ChatConversation[int(conversation.ID)].Member
	}

	client.Conn = conn
	client.IsConnected = true

	wsUseCase.Hub.Register <- client
	go client.Write()
	go client.Read(wsUseCase.Hub)
}

type ConversationInfo struct {
	Name     string `json:"name" binding:"required"`
	PosterID int    `json:"poster_id" binding:"required"`
}

// CreateConversation CreateOrCheckExistConversation godoc
// @Summary Create or check to exist a chat conversation for two users
// @Description Create or check to exist a chat conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation body ConversationInfo true "CreateConversation"
// @Success 200 {object} string
// @Router /chats/authorize/conversation [post]
func CreateConversation(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	var request ConversationInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())
	poster, err := chatRepo.GetPosterOwner(uint(request.PosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if poster.UserID == uint(payload.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't create conversation with yourself"})
		return
	}
	existConversation, err := chatRepo.ExistConversation(poster.UserID, uint(payload.UserID), poster.ID)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if existConversation.ID != 0 {
		c.JSON(http.StatusOK, gin.H{"message": "CreateConversation already exist", "conversation": existConversation})
		return
	}
	var conversationImage = ""
	if len(poster.Images) != 0 {
		conversationImage = poster.Images[0].Url
	}
	conversation, err := chatRepo.CreateConversation(request.Name, conversationImage, poster.UserID, uint(payload.UserID),
		uint(request.PosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CreateConversation created successfully", "conversation": conversation})
}

// AllUserConversations godoc
// @Summary Get all user conversations
// @Description Get all user conversations
// @Tags Chat
// @Accept json
// @Produce json
// @Success 200 {array} View.ConversationView
// @Router /chats/authorize/conversations [get]
func AllUserConversations(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())

	user, err := chatRepo.GetAllUserConversations(uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetAllUserConversation(c, user)
}

type GetConversationByIdPathRequest struct {
	ConversationID uint `uri:"conversation_id" binding:"required"`
}

// GetConversationById godoc
// @Summary Get conversation by id
// @Description Get conversation by id
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation_id path int true "Conversation ID"
// @Success 200 {object} string
// @Router /chats/authorize/conversation/{conversation_id} [get]
func GetConversationById(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	var pathRequest GetConversationByIdPathRequest
	if err := c.ShouldBindUri(&pathRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())
	conversation, err := chatRepo.GetUserConversationById(pathRequest.ConversationID, uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get conversation successfully", "conversation": conversation})
}

type ConversationHistoryPathRequest struct {
	ConversationID uint `uri:"conversation_id" binding:"required"`
}

type ConversationHistoryQueryRequest struct {
	PageID   int `form:"page_id" binding:"required"`
	PageSize int `form:"page_size" binding:"required,min=5"`
}

// ConversationHistory godoc
// @Summary Get conversation history
// @Description Get conversation history
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation_id path uint true "CreateConversation ID"
// @Param page_id query int true "Page ID" minimum(1) default(1)
// @Param page_size query int true "Page size" minimum(1) default(10)
// @Success 200 {array} Model.Conversation
// @Router /chats/authorize/history/{conversation_id}/ [get]
func ConversationHistory(c *gin.Context) {
	chatRepository := Repository.NewChatRepository(DBConfiguration.GetDB())

	var pathRequest ConversationHistoryPathRequest
	if err := c.ShouldBindUri(&pathRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var queryRequest ConversationHistoryQueryRequest
	if err := c.ShouldBindQuery(&queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset := (queryRequest.PageID - 1) * queryRequest.PageSize
	messages, err := chatRepository.GetConversationHistory(pathRequest.ConversationID, queryRequest.PageSize, offset)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetConversationHistory(c, messages)
}
