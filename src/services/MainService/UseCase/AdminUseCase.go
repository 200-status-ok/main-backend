package UseCase

import (
	"github.com/200-status-ok/main-backend/src/MainService/Cmd/DB"
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	Utils2 "github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/MainService/View"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type SignupAdminRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" binding:"required"`
	FName    string `json:"f_name"   binding:"required,min=4,max=30"`
	LName    string `json:"l_name"   binding:"required,min=4,max=30"`
	Email    string `json:"email"    binding:"required,min=8,max=30"`
	Phone    string `json:"phone"    binding:"required,min=11,max=30"`
}

func SignupAdminResponse(c *gin.Context) {
	db, _ := DB.GetDB()
	adminRepository := Repository.NewAdminRepository(db)
	var request SignupAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	admin := Model.Admin{
		Username: request.Username,
		Password: string(bcryptPassword),
		FName:    request.FName,
		LName:    request.LName,
		Email:    request.Email,
		Phone:    request.Phone,
	}
	admin, err = adminRepository.CreateAdmin(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	View.SignupAdminView(admin, c)
}

type LoginAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginAdminResponse(c *gin.Context) {
	db, _ := DB.GetDB()
	adminRepository := Repository.NewAdminRepository(db)
	var request LoginAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, err := adminRepository.GetAdminByUsername(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	jwtMaker, err := Token.NewJWTMaker(utils.ReadFromEnvFile(".env", "JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(request.Username, uint64(admin.ID), "Admin", time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.LoginAdminView(token, c)
}

type GetUserByAdminRequest struct {
	UserID uint `uri:"userid" binding:"required,min=1"`
}

func GetUserByIdAdminResponse(c *gin.Context) {
	var request GetUserByAdminRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, tx := DB.GetDB()
	userRepository := Repository.NewUserRepository(db, tx)
	user, err := userRepository.FindById(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUserByIdView(*user, c)
}

func GetUsersResponse(c *gin.Context) {
	db, tx := DB.GetDB()
	userRepository := Repository.NewUserRepository(db, tx)
	users, err := userRepository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUsersView(*users, c)
}

type UpdateUserByAdminRequest struct {
	UserID uint `uri:"userid" binding:"required,min=1"`
}

type UserInfo struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
}

func UpdateUserByIdAdminResponse(c *gin.Context) {
	var request UpdateUserByAdminRequest
	var userInfo UserInfo
	db, tx := DB.GetDB()
	userRepository := Repository.NewUserRepository(db, tx)
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils2.UsernameValidation(userInfo.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	updateUser, err := userRepository.UserUpdate(&Model.User{
		Username: userInfo.Username,
	}, request.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetUserByIdView(*updateUser, c)
}

type DeleteUserRequest struct {
	UserID uint `uri:"userid" binding:"required,min=1"`
}

func DeleteUserByIdAdminResponse(c *gin.Context) {
	var request DeleteUserRequest
	db, tx := DB.GetDB()
	userRepository := Repository.NewUserRepository(db, tx)
	esDeletePostersByUserId := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userRepository.DeleteUser(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = esDeletePostersByUserId.DeletePosterByUserID(int(request.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
