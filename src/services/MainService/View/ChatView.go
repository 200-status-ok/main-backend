package View

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
)

type ConversationView struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	IsOwner  bool   `json:"is_owner"`
	PosterID uint   `json:"poster_id"`
}

func GetAllUserConversation(c *gin.Context, user *Model.User) {
	var conversationView []ConversationView

	for _, conversation := range user.OwnConversations {
		conversationView = append(conversationView, ConversationView{
			ID:       conversation.ID,
			Name:     conversation.Name,
			ImageURL: conversation.ImageURL,
			IsOwner:  true,
			PosterID: conversation.PosterID,
		})
	}

	for _, conversation := range user.MemberConversations {
		conversationView = append(conversationView, ConversationView{
			ID:       conversation.ID,
			Name:     conversation.Name,
			ImageURL: conversation.ImageURL,
			IsOwner:  false,
			PosterID: conversation.PosterID,
		})
	}

	c.JSON(200, conversationView)
}

type ConversationHistoryView struct {
	Messages []Model.Message `json:"messages"`
	UserID   uint            `json:"user_id"`
}

func GetConversationHistory(c *gin.Context, messages []Model.Message, userID uint) {
	var conversationHistoryView ConversationHistoryView
	conversationHistoryView.Messages = messages
	conversationHistoryView.UserID = userID
	c.JSON(200, conversationHistoryView)
}
