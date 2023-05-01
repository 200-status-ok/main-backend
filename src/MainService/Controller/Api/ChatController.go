package Api

import (
	"github.com/gin-gonic/gin"
)

// JoinChat godoc
// @Summary JoinConversation a chat room
// @Description JoinConversation a chat room
// @Tags Chat
// @Accept json
// @Produce json
// @Param id query uint true "Chat ID"
// @Param client_id query uint true "Client ID"
// @Success 200 {object} string
// @Router /chats/join [get]
func JoinChat(c *gin.Context) {
	//UseCase.JoinConversation(c)
}
