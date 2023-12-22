package View

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConversationView struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	ImageURL    string          `json:"image_url"`
	IsOwner     bool            `json:"is_owner"`
	PosterID    uint            `json:"poster_id"`
	LastMessage ChatMessageView `json:"last_message"`
}

type ConversationByIDView struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	OwnerID   uint   `json:"owner_id"`
	MemberID  uint   `json:"member_id"`
	PosterID  uint   `json:"poster_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type ChatMessageView struct {
	ID             int64  `json:"id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	SenderID       uint   `json:"sender_id"`
	ReceiverID     uint   `json:"receiver_id"`
	Status         string `json:"status"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

func GetAllUserConversation(c *gin.Context, user *Model.User) {
	var conversationView []ConversationView
	for _, conversation := range user.OwnConversations {
		lastMessage := ChatMessageView{
			ID:             conversation.Messages[len(conversation.Messages)-1].ID,
			Content:        conversation.Messages[len(conversation.Messages)-1].Content,
			Type:           conversation.Messages[len(conversation.Messages)-1].Type,
			ConversationID: conversation.Messages[len(conversation.Messages)-1].ConversationId,
			SenderID:       conversation.Messages[len(conversation.Messages)-1].SenderId,
			ReceiverID:     conversation.Messages[len(conversation.Messages)-1].ReceiverId,
			Status:         conversation.Messages[len(conversation.Messages)-1].Status,
			CreatedAt:      conversation.Messages[len(conversation.Messages)-1].CreatedAt.Unix(),
			UpdatedAt:      conversation.Messages[len(conversation.Messages)-1].UpdatedAt.Unix(),
		}
		conversationView = append(conversationView, ConversationView{
			ID:          conversation.ID,
			Name:        conversation.Name,
			ImageURL:    conversation.ImageURL,
			IsOwner:     true,
			LastMessage: lastMessage,
			PosterID:    conversation.PosterID,
		})
	}

	for _, conversation := range user.MemberConversations {
		lastMessage := ChatMessageView{
			ID:             conversation.Messages[len(conversation.Messages)-1].ID,
			Content:        conversation.Messages[len(conversation.Messages)-1].Content,
			Type:           conversation.Messages[len(conversation.Messages)-1].Type,
			ConversationID: conversation.Messages[len(conversation.Messages)-1].ConversationId,
			SenderID:       conversation.Messages[len(conversation.Messages)-1].SenderId,
			ReceiverID:     conversation.Messages[len(conversation.Messages)-1].ReceiverId,
			Status:         conversation.Messages[len(conversation.Messages)-1].Status,
			CreatedAt:      conversation.Messages[len(conversation.Messages)-1].CreatedAt.Unix(),
			UpdatedAt:      conversation.Messages[len(conversation.Messages)-1].UpdatedAt.Unix(),
		}
		conversationView = append(conversationView, ConversationView{
			ID:          conversation.ID,
			Name:        conversation.Name,
			ImageURL:    conversation.ImageURL,
			IsOwner:     false,
			LastMessage: lastMessage,
			PosterID:    conversation.PosterID,
		})
	}

	c.JSON(200, conversationView)
}

type ConversationHistoryView struct {
	Messages []ChatMessageView `json:"messages"`
	UserID   uint              `json:"user_id"`
}

func GetConversationHistory(c *gin.Context, messages []Model.Message, userID uint) {
	var conversationHistoryView ConversationHistoryView
	messagesView := make([]ChatMessageView, 0)
	for _, message := range messages {
		messagesView = append(messagesView, ChatMessageView{
			ID:             message.ID,
			Content:        message.Content,
			Type:           message.Type,
			ConversationID: message.ConversationId,
			SenderID:       message.SenderId,
			ReceiverID:     message.ReceiverId,
			Status:         message.Status,
			CreatedAt:      message.CreatedAt.Unix(),
			UpdatedAt:      message.UpdatedAt.Unix(),
		})
	}
	conversationHistoryView.Messages = messagesView
	conversationHistoryView.UserID = userID
	c.JSON(200, conversationHistoryView)
}

func GetConversationByID(conversation Model.Conversation, c *gin.Context) {
	conversationByID := ConversationByIDView{
		ID:        conversation.ID,
		Name:      conversation.Name,
		ImageURL:  conversation.ImageURL,
		OwnerID:   conversation.OwnerID,
		MemberID:  conversation.MemberID,
		PosterID:  conversation.PosterID,
		CreatedAt: conversation.CreatedAt.Unix(),
		UpdatedAt: conversation.UpdatedAt.Unix(),
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get conversation successfully", "conversation": conversationByID})
}
