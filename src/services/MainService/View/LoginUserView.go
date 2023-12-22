package View

import "github.com/gin-gonic/gin"

type UserView struct {
	ID      uint   `json:"id"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type MessageView struct {
	Message string `json:"message"`
}

func LoginUserView(userID uint, token string, c *gin.Context) {
	result := UserView{
		ID:      userID,
		Token:   token,
		Message: "login successful",
	}
	c.JSON(200, result)
}

func LoginMessageView(message string, c *gin.Context) {
	result := MessageView{
		Message: message,
	}
	c.JSON(200, result)
}
