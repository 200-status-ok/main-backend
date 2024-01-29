package Api

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Cmd/DB"
	"github.com/200-status-ok/main-backend/src/MainService/RealtimeChat"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/200-status-ok/main-backend/src/MainService/View"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatWS struct {
	Hub *RealtimeChat.Hub
}

func NewChatWS(hub *RealtimeChat.Hub) *ChatWS {
	return &ChatWS{Hub: hub}
}

type MessageBody struct {
	ID             int64  `json:"id" binding:"required"`
	ConversationID int    `json:"conversation_id" binding:"required"`
	PosterID       uint   `json:"poster_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	Type           string `json:"type" binding:"required"`
}

// SendMessage godoc
// @Summary SendMessage
// @Description SendMessage to join a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Message body MessageBody true "Message"
// @Success 200 {object} string
// @Router /chat/authorize/message [post]
func (wsUseCase *ChatWS) SendMessage(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	db, _ := DB.GetDB()
	chatRepo := Repository.NewChatRepository(db)
	posterRepo := Repository.NewPosterRepository(db)
	var request MessageBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poster, err := posterRepo.GetPosterById(int(request.PosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	senderID := uint(payload.UserID)
	var receiverID uint
	if request.ConversationID == -1 {
		var conversationImage string
		receiverID = poster.UserID
		if receiverID == senderID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can't create conversation with yourself"})
			return
		}
		if len(poster.Images) == 0 {
			conversationImage = ""
		} else {
			conversationImage = poster.Images[0].Url
		}
		conversation, err := chatRepo.CreateConversation(poster.Title, conversationImage, receiverID, senderID,
			poster.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		request.ConversationID = int(conversation.ID)
	}
	existConv, err := chatRepo.ExistConversation(uint(request.ConversationID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if senderID == existConv.OwnerID {
		receiverID = existConv.MemberID
	} else {
		receiverID = existConv.OwnerID
	}
	if senderID == receiverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't send message to yourself"})
		return
	}
	sendTime := time.Now()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := chatRepo.SaveMessage(request.ID, uint(request.ConversationID), senderID, request.Content, request.Type,
		int(receiverID), sendTime, "unread")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg := &dtos.Message{
		ID:             message.ID,
		Content:        message.Content,
		ConversationID: int(message.ConversationId),
		SenderID:       int(message.SenderId),
		ReceiverId:     int(message.ReceiverId),
		SequenceNo:     message.SequenceNumber,
		Time:           message.CreatedAt.UnixMilli(),
		Type:           message.Type,
		Status:         message.Status,
	}
	wsUseCase.Hub.Broadcast <- msg
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "send_message": msg})
}

type OpenWSConnection struct {
	Token string `form:"token" binding:"required"`
}

// OpenWSConnection godoc
// @Summary OpenWSConnection
// @Description OpenWSConnection to join a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Message body dtos.TransferMessage true "Message"
// @Param token query string true "Token"
// @Success 200 {object} string
// @Router /chat/open-ws [get]
func (wsUseCase *ChatWS) OpenWSConnection(c *gin.Context) {
	db, _ := DB.GetDB()
	chatRepo := Repository.NewChatRepository(db)
	var request OpenWSConnection
	secretKey := utils.ReadFromEnvFile(".env", "JWT_SECRET")
	tokenMaker, _ := Token.NewJWTMaker(secretKey)
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	if _, ok := wsUseCase.Hub.Clients[int(payload.UserID)]; ok {
		if wsUseCase.Hub.Clients[int(payload.UserID)].Status == "online" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are already online"})
			fmt.Println("You are already online")
			return
		}
	}
	conversations, err := chatRepo.GetConversationByUserID(uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uniqueMap := make(map[int]struct{})
	pairedUsers := make([]int, 0)
	for _, conversation := range conversations {
		if conversation.OwnerID == uint(payload.UserID) {
			if _, ok := uniqueMap[int(conversation.MemberID)]; !ok {
				uniqueMap[int(conversation.MemberID)] = struct{}{}
				pairedUsers = append(pairedUsers, int(conversation.MemberID))
			}
		} else {
			if _, ok := uniqueMap[int(conversation.OwnerID)]; !ok {
				uniqueMap[int(conversation.OwnerID)] = struct{}{}
				pairedUsers = append(pairedUsers, int(conversation.OwnerID))
			}
		}
	}
	if _, ok := wsUseCase.Hub.Clients[int(payload.UserID)]; !ok {
		var client = RealtimeChat.Client{
			ID:      int(payload.UserID),
			Message: make(chan *dtos.Message, 100),
			Conn:    &websocket.Conn{},
			Status:  "offline",
		}
		wsUseCase.Hub.Clients[int(payload.UserID)] = &client
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	wsUseCase.Hub.PairUsers[int(payload.UserID)] = pairedUsers
	wsUseCase.Hub.Clients[int(payload.UserID)].Conn = conn
	wsUseCase.Hub.Clients[int(payload.UserID)].Status = "online"
	wsUseCase.Hub.Register <- wsUseCase.Hub.Clients[int(payload.UserID)]

	go wsUseCase.Hub.Clients[int(payload.UserID)].Write()
	go wsUseCase.Hub.Clients[int(payload.UserID)].UserTrace(wsUseCase.Hub)
}

type ConversationInfo struct {
	Name     string `json:"name" binding:"required"`
	PosterID int    `json:"poster_id" binding:"required"`
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
	db, _ := DB.GetDB()
	chatRepo := Repository.NewChatRepository(db)

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

	db, _ := DB.GetDB()
	chatRepo := Repository.NewChatRepository(db)
	conversation, err := chatRepo.GetUserConversationById(pathRequest.ConversationID, uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetConversationByID(conversation, c)
}

type ConversationIDPathRequest struct {
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
	db, _ := DB.GetDB()
	chatRepository := Repository.NewChatRepository(db)
	payload := c.MustGet("authorization_payload").(*Token.Payload)

	var pathRequest ConversationIDPathRequest
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

type UpdateConversationBody struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// UpdateConversation godoc
// @Summary Update conversation
// @Description Update conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation_id path uint true "CreateConversation ID"
// @Param UpdateConversation body UpdateConversationBody true "UpdateConversationBody"
// @Success 200 {object} string
// @Router /chat/authorize/conversation/{conversation_id} [patch]
func UpdateConversation(c *gin.Context) {
	db, _ := DB.GetDB()
	chatRepository := Repository.NewChatRepository(db)

	var pathRequest ConversationIDPathRequest
	if err := c.ShouldBindUri(&pathRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var bodyRequest UpdateConversationBody
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := chatRepository.UpdateConversation(pathRequest.ConversationID, bodyRequest.Name, bodyRequest.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation updated successfully"})
}

type MessageIDsBody struct {
	MessageIDs []int `json:"message_ids" binding:"required"`
	SenderID   int   `json:"sender_id" binding:"required"`
}

// ReadMessages godoc
// @Summary Read Multiple Messages
// @Description Read Multiple Messages
// @Tags Chat
// @Accept json
// @Produce json
// @Param MessageID body MessageIDsBody true "MessageIDs"
// @Success 200 {object} string
// @Router /chat/authorize/read [post]
func (wsUseCase *ChatWS) ReadMessages(c *gin.Context) {
	db, _ := DB.GetDB()
	chatRepository := Repository.NewChatRepository(db)
	var request MessageIDsBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := chatRepository.ReadMessages(request.MessageIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content := fmt.Sprintf("%v", request.MessageIDs)
	messageTime := time.Now()
	msg := &dtos.Message{
		ID:             0,
		Content:        content,
		ConversationID: 0,
		SenderID:       0,
		ReceiverId:     request.SenderID,
		Time:           messageTime.UnixMilli(),
		Type:           "text-read-message",
		Status:         "notification",
	}
	wsUseCase.Hub.Broadcast <- msg
	c.JSON(http.StatusOK, gin.H{"message": "Messages read successfully"})
}
