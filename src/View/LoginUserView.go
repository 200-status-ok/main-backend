package View

import "github.com/gin-gonic/gin"

type UserView struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func LoginUserView(token string, c *gin.Context) {
	result := UserView{
		Token:   token,
		Message: "Login successful",
	}
	c.JSON(200, result)
}
