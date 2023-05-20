package View

import (
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserViewPayments struct {
	Id        uint    `json:"id"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
	Status    string  `json:"status"`
	UserID    uint    `json:"user_id"`
}
type UserViewID struct {
	Id           uint                  `json:"id"`
	Username     string                `json:"username"`
	Wallet       float64               `json:"wallet"`
	Posters      []Model2.Poster       `json:"posters"`
	MarkedPoster []Model2.MarkedPoster `json:"marked_posters"`
}
type UserViewIDs struct {
	Id           uint                  `json:"id"`
	Username     string                `json:"username"`
	Wallet       float64               `json:"wallet"`
	Posters      []Model2.Poster       `json:"posters"`
	MarkedPoster []Model2.MarkedPoster `json:"marked_posters"`
}

func GetUserByIdView(user Model2.User, c *gin.Context) {
	result := UserViewID{
		Id:           user.ID,
		Username:     user.Username,
		Posters:      user.Posters,
		MarkedPoster: user.MarkedPosters,
		Wallet:       user.Wallet,
	}
	c.JSON(http.StatusOK, result)
}

func GetUsersView(users []Model2.User, c *gin.Context) {
	var result []UserViewIDs
	for _, user := range users {
		result = append(result, UserViewIDs{
			Id:           user.ID,
			Username:     user.Username,
			Wallet:       user.Wallet,
			Posters:      user.Posters,
			MarkedPoster: user.MarkedPosters,
		})
	}
	c.JSON(http.StatusOK, result)
}

func GetUserPaymentsView(payments []Model2.Payment, c *gin.Context) {
	var result []UserViewPayments
	for _, payment := range payments {
		result = append(result, UserViewPayments{
			Id:        payment.ID,
			Amount:    payment.Amount,
			CreatedAt: payment.CreatedAt.String(),
			Status:    payment.Status,
			UserID:    payment.UserID,
		})
	}
	c.JSON(http.StatusOK, result)
}
