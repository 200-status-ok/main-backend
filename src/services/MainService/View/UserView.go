package View

import (
	Model2 "github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserViewPayment struct {
	Id        uint    `json:"id"`
	Amount    float64 `json:"amount"`
	CreatedAt int64   `json:"created_at"`
	Status    string  `json:"status"`
	UserID    uint    `json:"user_id"`
}

type UserViewInfo struct {
	Id           uint               `json:"id"`
	Username     string             `json:"username"`
	Wallet       float64            `json:"wallet"`
	Posters      []PosterView       `json:"posters"`
	MarkedPoster []MarkedPosterView `json:"marked_posters"`
}

func GetUserByIdView(user Model2.User, c *gin.Context) {
	var userPosters []PosterView
	for _, poster := range user.Posters {
		addressesView := make([]AddressView, 0)
		for _, address := range poster.Addresses {
			addressesView = append(addressesView, AddressView{
				ID:            address.ID,
				Province:      address.Province,
				City:          address.City,
				AddressDetail: address.AddressDetail,
				Latitude:      address.Latitude,
				Longitude:     address.Longitude,
				PosterID:      address.PosterId,
				CreatedAt:     address.CreatedAt.UnixMilli(),
				UpdatedAt:     address.UpdatedAt.UnixMilli(),
			})
		}
		imagesView := make([]ImageView, 0)
		for _, image := range poster.Images {
			imagesView = append(imagesView, ImageView{
				ID:        image.ID,
				Url:       image.Url,
				CreatedAt: image.CreatedAt.UnixMilli(),
				UpdatedAt: image.UpdatedAt.UnixMilli(),
			})
		}
		tagsView := make([]TagView, 0)
		for _, tag := range poster.Tags {
			tagsView = append(tagsView, TagView{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt.UnixMilli(),
				UpdatedAt: tag.UpdatedAt.UnixMilli(),
			})
		}
		userPosters = append(userPosters, PosterView{
			ID:          poster.ID,
			Title:       poster.Title,
			Description: poster.Description,
			Addresses:   addressesView,
			Images:      imagesView,
			TelegramId:  poster.TelegramID,
			UserPhone:   poster.UserPhone,
			Status:      poster.Status,
			Tags:        tagsView,
			User:        poster.UserID,
			Award:       poster.Award,
			CreatedAt:   poster.CreatedAt.UnixMilli(),
			UpdatedAt:   poster.UpdatedAt.UnixMilli(),
			State:       poster.State,
		})
	}
	markedPosters := make([]MarkedPosterView, 0)
	for _, markedPoster := range user.MarkedPosters {
		addressesView := make([]AddressView, 0)
		for _, address := range markedPoster.Poster.Addresses {
			addressesView = append(addressesView, AddressView{
				ID:            address.ID,
				Province:      address.Province,
				City:          address.City,
				AddressDetail: address.AddressDetail,
				Latitude:      address.Latitude,
				Longitude:     address.Longitude,
				PosterID:      address.PosterId,
				CreatedAt:     address.CreatedAt.UnixMilli(),
				UpdatedAt:     address.UpdatedAt.UnixMilli(),
			})
		}
		imagesView := make([]ImageView, 0)
		for _, image := range markedPoster.Poster.Images {
			imagesView = append(imagesView, ImageView{
				ID:        image.ID,
				Url:       image.Url,
				CreatedAt: image.CreatedAt.UnixMilli(),
				UpdatedAt: image.UpdatedAt.UnixMilli(),
			})
		}
		tagsView := make([]TagView, 0)
		for _, tag := range markedPoster.Poster.Tags {
			tagsView = append(tagsView, TagView{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt.UnixMilli(),
				UpdatedAt: tag.UpdatedAt.UnixMilli(),
			})
		}
		poster := PosterView{
			ID:          markedPoster.Poster.ID,
			Title:       markedPoster.Poster.Title,
			Description: markedPoster.Poster.Description,
			Addresses:   addressesView,
			Images:      imagesView,
			TelegramId:  markedPoster.Poster.TelegramID,
			UserPhone:   markedPoster.Poster.UserPhone,
			Status:      markedPoster.Poster.Status,
			Tags:        tagsView,
			User:        markedPoster.Poster.UserID,
			Award:       markedPoster.Poster.Award,
			CreatedAt:   markedPoster.Poster.CreatedAt.UnixMilli(),
			UpdatedAt:   markedPoster.Poster.UpdatedAt.UnixMilli(),
			State:       markedPoster.Poster.State,
		}
		markedPosters = append(markedPosters, MarkedPosterView{
			ID:        markedPoster.ID,
			PosterID:  markedPoster.PosterID,
			UserID:    markedPoster.UserID,
			Poster:    poster,
			CreatedAt: markedPoster.CreatedAt.UnixMilli(),
			UpdatedAt: markedPoster.UpdatedAt.UnixMilli(),
		})
	}
	result := UserViewInfo{
		Id:           user.ID,
		Username:     user.Username,
		Posters:      userPosters,
		MarkedPoster: markedPosters,
		Wallet:       user.Wallet,
	}
	c.JSON(http.StatusOK, result)
}

func GetUsersView(users []Model2.User, c *gin.Context) {
	var result []UserViewInfo
	for _, user := range users {
		var userPosters []PosterView
		markedPosters := make([]MarkedPosterView, 0)
		for _, markedPoster := range user.MarkedPosters {
			addressesView := make([]AddressView, 0)
			for _, address := range markedPoster.Poster.Addresses {
				addressesView = append(addressesView, AddressView{
					ID:            address.ID,
					Province:      address.Province,
					City:          address.City,
					AddressDetail: address.AddressDetail,
					Latitude:      address.Latitude,
					Longitude:     address.Longitude,
					PosterID:      address.PosterId,
					CreatedAt:     address.CreatedAt.UnixMilli(),
					UpdatedAt:     address.UpdatedAt.UnixMilli(),
				})
			}
			imagesView := make([]ImageView, 0)
			for _, image := range markedPoster.Poster.Images {
				imagesView = append(imagesView, ImageView{
					ID:        image.ID,
					Url:       image.Url,
					CreatedAt: image.CreatedAt.UnixMilli(),
					UpdatedAt: image.UpdatedAt.UnixMilli(),
				})
			}
			tagsView := make([]TagView, 0)
			for _, tag := range markedPoster.Poster.Tags {
				tagsView = append(tagsView, TagView{
					ID:        tag.ID,
					Name:      tag.Name,
					CreatedAt: tag.CreatedAt.UnixMilli(),
					UpdatedAt: tag.UpdatedAt.UnixMilli(),
				})
			}
			poster := PosterView{
				ID:          markedPoster.Poster.ID,
				Title:       markedPoster.Poster.Title,
				Description: markedPoster.Poster.Description,
				Addresses:   addressesView,
				Images:      imagesView,
				TelegramId:  markedPoster.Poster.TelegramID,
				UserPhone:   markedPoster.Poster.UserPhone,
				Status:      markedPoster.Poster.Status,
				Tags:        tagsView,
				User:        markedPoster.Poster.UserID,
				Award:       markedPoster.Poster.Award,
				CreatedAt:   markedPoster.Poster.CreatedAt.UnixMilli(),
				UpdatedAt:   markedPoster.Poster.UpdatedAt.UnixMilli(),
				State:       markedPoster.Poster.State,
			}
			markedPosters = append(markedPosters, MarkedPosterView{
				ID:        markedPoster.ID,
				PosterID:  markedPoster.PosterID,
				UserID:    markedPoster.UserID,
				Poster:    poster,
				CreatedAt: markedPoster.CreatedAt.UnixMilli(),
				UpdatedAt: markedPoster.UpdatedAt.UnixMilli(),
			})
		}
		for _, poster := range user.Posters {
			addressesView := make([]AddressView, 0)
			for _, address := range poster.Addresses {
				addressesView = append(addressesView, AddressView{
					ID:            address.ID,
					Province:      address.Province,
					City:          address.City,
					AddressDetail: address.AddressDetail,
					Latitude:      address.Latitude,
					Longitude:     address.Longitude,
					PosterID:      address.PosterId,
					CreatedAt:     address.CreatedAt.UnixMilli(),
					UpdatedAt:     address.UpdatedAt.UnixMilli(),
				})
			}
			imagesView := make([]ImageView, 0)
			for _, image := range poster.Images {
				imagesView = append(imagesView, ImageView{
					ID:        image.ID,
					Url:       image.Url,
					CreatedAt: image.CreatedAt.UnixMilli(),
					UpdatedAt: image.UpdatedAt.UnixMilli(),
				})
			}
			tagsView := make([]TagView, 0)
			for _, tag := range poster.Tags {
				tagsView = append(tagsView, TagView{
					ID:        tag.ID,
					Name:      tag.Name,
					CreatedAt: tag.CreatedAt.UnixMilli(),
					UpdatedAt: tag.UpdatedAt.UnixMilli(),
				})
			}
			userPosters = append(userPosters, PosterView{
				ID:          poster.ID,
				Title:       poster.Title,
				Description: poster.Description,
				Addresses:   addressesView,
				Images:      imagesView,
				TelegramId:  poster.TelegramID,
				UserPhone:   poster.UserPhone,
				Status:      poster.Status,
				Tags:        tagsView,
				User:        poster.UserID,
				Award:       poster.Award,
				CreatedAt:   poster.CreatedAt.UnixMilli(),
				UpdatedAt:   poster.UpdatedAt.UnixMilli(),
				State:       poster.State,
			})
		}
		result = append(result, UserViewInfo{
			Id:           user.ID,
			Username:     user.Username,
			Wallet:       user.Wallet,
			Posters:      userPosters,
			MarkedPoster: markedPosters,
		})
	}
	c.JSON(http.StatusOK, result)
}

func GetUserPaymentsView(payments []Model2.Payment, c *gin.Context) {
	var result []UserViewPayment
	for _, payment := range payments {
		result = append(result, UserViewPayment{
			Id:        payment.ID,
			Amount:    payment.Amount,
			CreatedAt: payment.CreatedAt.UnixMilli(),
			Status:    payment.Status,
			UserID:    payment.UserID,
		})
	}
	c.JSON(http.StatusOK, result)
}
