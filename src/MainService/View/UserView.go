package View

import (
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserViewID struct {
	Id           uint                  `json:"id"`
	Username     string                `json:"username"`
	Posters      []Model2.Poster       `json:"posters"`
	MarkedPoster []Model2.MarkedPoster `json:"marked_posters"`
}

func GetUserByIdView(user Model2.User, c *gin.Context) {
	result := UserViewID{
		Id:           user.ID,
		Username:     user.Username,
		Posters:      user.Posters,
		MarkedPoster: user.MarkedPosters,
	}
	c.JSON(http.StatusOK, result)
}

func GetUsersView(users []Model2.User, c *gin.Context) {
	var result []UserViewID
	for _, user := range users {
		result = append(result, UserViewID{
			Id:           user.ID,
			Username:     user.Username,
			Posters:      user.Posters,
			MarkedPoster: user.MarkedPosters,
		})
	}
	c.JSON(http.StatusOK, result)
}
