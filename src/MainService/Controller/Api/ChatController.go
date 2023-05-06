package Api

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
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

type JoinConversationReq struct {
	ConversationID uint `form:"conv_id" binding:"required"`
	ClientID       uint `form:"client_id" binding:"required"`
}

type CreateConversation struct {
	Name     string `json:"name" binding:"required"`
	ClientID int    `json:"client_id" binding:"required"`
	PosterID int    `json:"poster_id" binding:"required"`
}

type ChatWS struct {
	Hub *WebSocket.Hub
}

func NewChatWS(hub *WebSocket.Hub) *ChatWS {
	return &ChatWS{Hub: hub}
}

// JoinConversation JoinChat godoc
// @Summary JoinConversation a chat room
// @Description JoinConversation a chat room
// @Tags Chat
// @Accept json
// @Produce json
// @Param id query uint true "Chat ID"
// @Param client_id query uint true "Client ID"
// @Success 200 {object} string
// @Router /chats/join [get]
func (wsUseCase *ChatWS) JoinConversation(c *gin.Context) {
	var request JoinConversationReq
	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	if request.ClientID == conversation.OwnerID {
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

// CreateChatConversation CreateConversation godoc
// @Summary Create a chat conversation for two users
// @Description Create a chat conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param room body CreateConversation true "ChatConversation"
// @Success 200 {object} string
// @Router /chats/conversation [post]
func CreateChatConversation(c *gin.Context) {
	var request CreateConversation
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
	err = chatRepo.CreateConversation(request.Name, poster.UserID, uint(request.ClientID), uint(request.PosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation created successfully"})
}
