package UseCase

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/WebSocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatWS struct {
	Hub *WebSocket.Hub
}

func NewChatWS(hub *WebSocket.Hub) *ChatWS {
	return &ChatWS{Hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type JoinConversationRequest struct {
	PosterID uint `form:"id" binding:"required"`
	ClientID uint `form:"client_id" binding:"required"`
}

type CreateConversation struct {
	Name     string `json:"name" binding:"required"`
	ClientID int    `json:"client_id" binding:"required"`
	PosterID int    `json:"poster_id" binding:"required"`
}

// CreateChatConversation CreateRoom Create CreateConversation godoc
// @Summary Create a chat conversation for two users
// @Description Create a chat conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param room body CreateConversation true "ChatConversation"
// @Success 200 {object} string
// @Router /chats/conversation [post]
func (wsUseCase *ChatWS) CreateChatConversation(c *gin.Context) {
	var request CreateConversation
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())
	_, err := chatRepo.GetChatRoomByPosterId(uint(request.PosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation created successfully"})
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
	var request JoinConversationRequest
	chatRepo := Repository.NewChatRepository(DBConfiguration.GetDB())

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chatRoom, err := chatRepo.GetChatRoomByPosterId(request.PosterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	conversation, err := chatRepo.GetConversationById(chatRoom.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if _, ok := wsUseCase.Hub.ChatConversation[int(conversation.ID)]; !ok {
		wsUseCase.Hub.ChatConversation[int(conversation.ID)] = &WebSocket.ConversationChat{
			ID:      int(conversation.ID),
			Name:    fmt.Sprintf("conversation-%d", conversation.ID),
			Members: map[int]*WebSocket.Client{},
			ChatRoom: &WebSocket.ChatRoom{
				ID:   int(chatRoom.ID),
				Name: fmt.Sprintf("chat-room-%d", chatRoom.ID),
			},
		}
	}

	client := &WebSocket.Client{
		Conn:           conn,
		Message:        make(chan *WebSocket.Message, 10),
		ID:             int(request.ClientID),
		Role:           WebSocket.Owner,
		ConversationID: int(conversation.ID),
	}

	if request.ClientID == chatRoom.OwnerID {
		client.Role = WebSocket.Owner
		wsUseCase.Hub.ChatConversation[int(conversation.ID)].Members[int(client.ID)] = client
		wsUseCase.Hub.Register <- client
		go client.Write()
		client.Read(wsUseCase.Hub)
	} else {
		client.Role = WebSocket.Member
		wsUseCase.Hub.ChatConversation[int(conversation.ID)].Members[client.ID] = client
		message := &WebSocket.Message{
			Content:        fmt.Sprintf("A new user {%d} has joined the chat", client.ID),
			ConversationID: int(conversation.ID),
			SenderID:       client.ID,
		}
		wsUseCase.Hub.Register <- client
		wsUseCase.Hub.Broadcast <- message
		go client.Write()
		client.Read(wsUseCase.Hub)
	}
}
