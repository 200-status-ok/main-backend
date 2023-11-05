package Api

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/RealtimeChat"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/200-status-ok/main-backend/src/MainService/View"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
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

type ChatWS2 struct {
	Hub *RealtimeChat.Hub
}

func NewChatWS(hub *RealtimeChat.Hub) *ChatWS2 {
	return &ChatWS2{Hub: hub}
}

type JoinConversationReq struct {
	Token string `form:"token" binding:"required"`
}

// OpenWSConnection godoc
// @Summary OpenWSConnection
// @Description OpenWSConnection to join a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Message body RealtimeChat.TransferMessage true "Message"
// @Param token query string true "Token"
// @Success 200 {object} string
// @Router /chat/open-ws [get]
func (wsUseCase *ChatWS2) OpenWSConnection(c *gin.Context) {
	chatRepo := Repository.NewChatRepository(pgsql.GetDB())
	var request JoinConversationReq
	secretKey := utils.ReadFromEnvFile(".env", "JWT_SECRET")
	tokenMaker, _ := Token.NewJWTMaker(secretKey)
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := request.Token
	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {
		fmt.Println("Token is invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conversations, err := chatRepo.GetConversationByUserID(uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pairedUsers := make([]int, 0)
	for _, conversation := range conversations {
		if conversation.OwnerID == uint(payload.UserID) {
			pairedUsers = append(pairedUsers, int(conversation.MemberID))
		} else {
			pairedUsers = append(pairedUsers, int(conversation.OwnerID))
		}
	}

	for _, conversation := range conversations {
		if _, ok := wsUseCase.Hub.Clients[int(conversation.OwnerID)]; !ok {
			var client = RealtimeChat.Client{
				ID:      int(conversation.OwnerID),
				Message: make(chan *dtos.Message, 100),
				Conn:    &websocket.Conn{},
			}
			wsUseCase.Hub.Clients[int(conversation.OwnerID)] = &client
		}
		if _, ok := wsUseCase.Hub.Clients[int(conversation.MemberID)]; !ok {
			var client = RealtimeChat.Client{
				ID:      int(conversation.MemberID),
				Message: make(chan *dtos.Message, 100),
				Conn:    &websocket.Conn{},
			}
			wsUseCase.Hub.Clients[int(conversation.MemberID)] = &client
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	wsUseCase.Hub.PairUsers[int(payload.UserID)] = pairedUsers
	wsUseCase.Hub.Clients[int(payload.UserID)].Conn = conn
	wsUseCase.Hub.Register <- wsUseCase.Hub.Clients[int(payload.UserID)]

	fmt.Println(wsUseCase.Hub.PairUsers[int(payload.UserID)])
	go wsUseCase.Hub.Clients[int(payload.UserID)].Write()
	go wsUseCase.Hub.Clients[int(payload.UserID)].Read(wsUseCase.Hub)
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
// @Router /chat/authorize/conversation [post]
func CreateConversation(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	var request ConversationInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chatRepo := Repository.NewChatRepository(pgsql.GetDB())
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
// @Router /chat/authorize/conversation [get]
func AllUserConversations(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	chatRepo := Repository.NewChatRepository(pgsql.GetDB())

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
// @Router /chat/authorize/conversation/{conversation_id} [get]
func GetConversationById(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	var pathRequest GetConversationByIdPathRequest
	if err := c.ShouldBindUri(&pathRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatRepo := Repository.NewChatRepository(pgsql.GetDB())
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
// @Router /chat/authorize/history/{conversation_id}/ [get]
func ConversationHistory(c *gin.Context) {
	chatRepository := Repository.NewChatRepository(pgsql.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)

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

	View.GetConversationHistory(c, messages, uint(payload.UserID))
}

// UpdateConversation godoc
// @Summary Update conversation
// @Description Update conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation_id path uint true "CreateConversation ID"
// @Param name body string true "Name"
// @Param image body string true "Image"
// @Success 200 {object} string
// @Router /chat/authorize/conversation/{conversation_id} [patch]
func UpdateConversation(c *gin.Context) {
}

// ReadMessageInConversation godoc
// @Summary Read conversation
// @Description Read conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation_id path uint true "CreateConversation ID"
// @Success 200 {object} string
// @Router /chat/authorize/read/{conversation_id} [get]
func ReadMessageInConversation(c *gin.Context) {
	chatRepository := Repository.NewChatRepository(pgsql.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)

	var pathRequest ConversationHistoryPathRequest
	if err := c.ShouldBindUri(&pathRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMessage, err := chatRepository.ReadMessageInConversation(pathRequest.ConversationID, uint(payload.UserID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.ReadMessageInConversationView(c, updatedMessage)
}
