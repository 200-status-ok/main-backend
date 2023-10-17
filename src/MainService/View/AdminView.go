package View

import (
	Model2 "github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
)

type AdminView struct {
	Username string `json:"username"`
	FName    string `json:"f_name"`
	LName    string `json:"l_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type AdminLoginView struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func LoginAdminView(token string, c *gin.Context) {
	result := UserView{
		Token:   token,
		Message: "Login successful",
	}
	c.JSON(200, result)
}

func SignupAdminView(admin Model2.Admin, c *gin.Context) {
	result := AdminView{
		Username: admin.Username,
		FName:    admin.FName,
		LName:    admin.LName,
		Email:    admin.Email,
		Phone:    admin.Phone,
	}
	c.JSON(200, result)
}
