package UseCase

import (
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/MainService/View"
	DTO2 "github.com/200-status-ok/main-backend/src/MainService/dtos"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GetPostersRequest struct {
	PageID       int     `form:"page_id" binding:"required,min=1" json:"page_id"`
	PageSize     int     `form:"page_size" binding:"required,min=5" json:"page_size"`
	Sort         string  `form:"sort,omitempty" binding:"omitempty,oneof=asc desc" json:"sort"`
	SortBy       string  `form:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at id" json:"sort_by"`
	Status       string  `form:"status,omitempty" binding:"oneof=lost found both" json:"status"`
	SearchPhrase string  `form:"search_phrase,omitempty" json:"search_phrase"`
	TimeStart    int64   `form:"time_start,omitempty" json:"time_start"`
	TimeEnd      int64   `form:"time_end,omitempty" json:"time_end"`
	OnlyAwards   bool    `form:"only_awards,omitempty" json:"only_awards"`
	Lat          float64 `form:"lat,omitempty" json:"lat"`
	Lon          float64 `form:"lon,omitempty" json:"lon"`
	TagIds       []int   `form:"tag_ids,omitempty" swaggertype:"array,int" json:"tag_ids"`
	State        string  `form:"state,omitempty" binding:"omitempty" json:"state"`
	SpecialType  string  `form:"special_type,omitempty" binding:"omitempty,oneof=all normal premium" json:"special_type"`
}

func GetPostersResponse(c *gin.Context) {
	esPosterCli := ElasticSearch.NewPosterES(elasticsearch.GetElastic())

	var request GetPostersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Sort = c.DefaultQuery("sort", "desc")
	request.SortBy = c.DefaultQuery("sort_by", "created_at")
	request.Status = c.DefaultQuery("status", "both")
	request.SearchPhrase = c.DefaultQuery("search_phrase", "")
	request.State = c.DefaultQuery("state", "all")
	request.SpecialType = c.DefaultQuery("special_type", "all")

	filterObject := DTO2.FilterObject{
		PageSize:     request.PageSize,
		Offset:       (request.PageID - 1) * request.PageSize,
		Sort:         request.Sort,
		SortBy:       request.SortBy,
		Status:       request.Status,
		SearchPhrase: request.SearchPhrase,
		TimeStart:    request.TimeStart,
		TimeEnd:      request.TimeEnd,
		OnlyAwards:   request.OnlyAwards,
		Lat:          request.Lat,
		Lon:          request.Lon,
		TagIds:       request.TagIds,
		State:        request.State,
		SpecialType:  request.SpecialType,
	}

	posters, total, err := esPosterCli.GetPosters(filterObject)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetPostersView(posters, total, request.PageSize, c)
}

type GetPosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func GetPosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
	var request GetPosterByIdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poster, err := posterRepository.GetPosterById(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetPosterByIdView(poster, c)
}

type DeletePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func DeletePosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)

	var request DeletePosterByIdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := posterRepository.DeletePosterById(uint(request.ID), uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poster deleted"})
}

type CreatePosterRequest struct {
	Poster    DTO2.CreatePosterDTO
	Addresses []DTO2.CreateAddressDTO
	ImgUrls   []string `json:"img_urls" binding:"required" json:"img_urls"`
	Tags      []string `json:"tags" binding:"required" json:"tags"`
}

func CreatePosterResponse(c *gin.Context) {
	var specialAdsPrice = 100000.0
	esPosterCli := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	var request CreatePosterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletAmount, err := userRepository.GetAmount(uint(payload.UserID))
	if request.Poster.SpecialType == "premium" {
		if walletAmount >= specialAdsPrice {
			_, err = userRepository.UpdateWallet(uint(payload.UserID), -specialAdsPrice)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			poster, err2 := posterRepository.CreatePoster(payload.UserID, request.Poster, request.Addresses, request.ImgUrls, request.Tags, "premium")
			if err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
				return
			}

			SendToNSFWQueue(poster.ID, c)
			TagValidationQueue(request.Tags, c)
			View.CreatePosterView(poster, c)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not enough money"})
			return
		}
	} else {
		poster, err := posterRepository.CreatePoster(payload.UserID, request.Poster, request.Addresses, request.ImgUrls, request.Tags, "normal")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		go func() {
			if request.Poster.Status == "found" {
				similarPosters, err := esPosterCli.FindSimilarityPosters(poster)
				if err != nil {
					fmt.Println(err)
				} else if len(similarPosters) > 0 {
					appEnv := os.Getenv("APP_ENV2")
					messageBroker := Utils.MessageClient{}
					if appEnv == "development" {
						err = messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
						if err != nil {
							c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
							return
						}
					} else if appEnv == "production" {
						err = messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
						if err != nil {
							c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
							return
						}
					}
					user, err := userRepository.FindById(similarPosters[0].UserID)
					if err != nil {
						fmt.Println(err)
					}
					if Utils.UsernameValidation(user.Username) == 0 {
						msg := "email/similar-poster/" + similarPosters[0].Title + "/" + user.Username
						err = messageBroker.PublishOnQueue([]byte(msg), "email_notification")
						if err != nil {
							fmt.Println(err)
						}
					} else {
						msg := "sms/similar-poster/" + similarPosters[0].Title + "/" + user.Username
						err = messageBroker.PublishOnQueue([]byte(msg), "sms_notification")
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}()
		SendToNSFWQueue(poster.ID, c)
		TagValidationQueue(request.Tags, c)
		View.CreatePosterView(poster, c)
	}
}

type UpdatePosterRequest struct {
	Poster    DTO2.UpdatePosterDTO    `json:"poster"`
	Addresses []DTO2.UpdateAddressDTO `json:"addresses"`
}

type UpdatePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func UpdatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)

	var request UpdatePosterRequest
	var posterByIdRequest UpdatePosterByIdRequest
	if err := c.ShouldBindUri(&posterByIdRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Poster.UserID = uint(payload.UserID)
	err := posterRepository.UpdatePoster(posterByIdRequest.ID, payload.Role, request.Poster, request.Addresses)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(payload.Role) == "user" {
		SendToNSFWQueue(uint(posterByIdRequest.ID), c)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poster updated successfully"})
}

func TagValidationQueue(tags []string, c *gin.Context) {
	appEnv := os.Getenv("APP_ENV2")
	messageBroker := Utils.MessageClient{}

	if appEnv == "development" {
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if appEnv == "production" {
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	msg := strings.Join(tags, ",")
	err := messageBroker.PublishOnQueue([]byte(msg), "tag-validation")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageBroker.Close()
}

func SendToNSFWQueue(posterID uint, c *gin.Context) {
	appEnv := os.Getenv("APP_ENV2")
	messageBroker := Utils.MessageClient{}

	if appEnv == "development" {
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if appEnv == "production" {
		err := messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	msg := strconv.Itoa(int(posterID))
	err := messageBroker.PublishOnQueue([]byte(msg), "nsfw-validation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageBroker.Close()
}

type CreatePosterReportRequest struct {
	PosterId    uint   `form:"poster_id" binding:"required"`
	IssuerId    uint   `form:"issuer_id" binding:"required"`
	ReportType  string `form:"report_type" binding:"required,oneof=spam inappropriate other"` //TODO: add more report types
	Description string `form:"description"`
}

func CreatePosterReportResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
	var request CreatePosterReportRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := posterRepository.CreatePosterReport(request.PosterId, request.IssuerId, request.ReportType, request.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Report created"})
}

type getPosterReportsRequest struct {
	PageID   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=20"`
	Status   string `form:"status,omitempty" binding:"oneof=open resolved both"`
}

func GetPosterReportsResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())

	var request getPosterReportsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset := (request.PageID - 1) * request.PageSize
	request.Status = c.DefaultQuery("status", "both")

	posterReports, err := posterRepository.GetAllPosterReports(request.PageSize, offset, request.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	View.GetPosterReportsView(posterReports, c)
}

type GetPosterReportByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func GetPosterReportByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())

	var request GetPosterReportByIdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posterReport, err := posterRepository.GetPosterReportById(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//DBConfiguration.CloseDB()
	View.GetPosterReportByIdView(posterReport, c)
}

type UpdatePosterReportRequest struct {
	PosterID    uint   `json:"poster_id"`
	IssuerID    uint   `json:"issuer_id"`
	ReportType  string `json:"report_type"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdatePosterReportByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func UpdatePosterReportResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())

	var request UpdatePosterReportRequest
	var id UpdatePosterReportByIdRequest

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := posterRepository.UpdatePosterReport(id.ID, request.PosterID, request.IssuerID, request.ReportType, request.Description, request.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Report resolved"})
}

type UpdatePosterStateRequest struct {
	ID    uint   `form:"id" binding:"required,min=1"`
	State string `form:"state" binding:"required,oneof=accepted rejected pending"`
}

func UpdatePosterStateResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(pgsql.GetDB())

	var request UpdatePosterStateRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := posterRepository.UpdatePosterState(request.ID, request.State)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poster state updated!"})
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Tag struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	DisplayText string `json:"display_text"`
}

type District struct {
	Level      string     `json:"level"`
	Radius     int        `json:"radius"`
	Centroid   Coordinate `json:"centroid"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	New        bool       `json:"new"`
	Parent     int        `json:"parent"`
	Neighbors  []int      `json:"neighbors"`
	SecondSlug string     `json:"second_slug"`
	Tags       []Tag      `json:"tags"`
}

type Data struct {
	Districts []District `json:"districts"`
}

func MockPoster(count int, userID int, tagNames []string) error {
	file, err := os.Open("Utils/Tehran.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
	var data Data
	err = json.Unmarshal(content, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	file2, err := ioutil.ReadFile("Utils/ObjectsFA.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	newString := strings.ReplaceAll(string(file2), "\r", "")
	splitStr := strings.Split(newString, "\n")
	// i 0 to count
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		randomNumber := rand.Intn(len(data.Districts))
		fmt.Println(data.Districts[randomNumber].Name)
		randomNumber2 := rand.Intn(len(splitStr))
		fmt.Println(splitStr[randomNumber2])
		randomNumber3 := rand.Intn(2)
		fmt.Println(randomNumber3)
		randomNumber4 := rand.Intn(len(tagNames))
		if randomNumber3 == 0 {
			request := CreatePosterRequest{
				Poster: DTO2.CreatePosterDTO{
					Title:       splitStr[randomNumber2],
					Description: "من یک " + splitStr[randomNumber2] + " گم کردم",
					Status:      "lost",
					Alert:       true,
					Chat:        true,
				},
				Addresses: []DTO2.CreateAddressDTO{
					{
						Province:      "تهران",
						City:          "تهران",
						AddressDetail: data.Districts[randomNumber].Name,
						Latitude:      data.Districts[randomNumber].Centroid.Latitude,
						Longitude:     data.Districts[randomNumber].Centroid.Longitude,
					},
				},
				ImgUrls: []string{},
				Tags: []string{
					tagNames[randomNumber4],
				},
			}
			posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
			model, err := posterRepository.CreatePoster(uint64(userID), request.Poster, request.Addresses, nil,
				request.Tags, "normal")
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println(model)
		} else {
			request := CreatePosterRequest{
				Poster: DTO2.CreatePosterDTO{
					Title:       splitStr[randomNumber2],
					Description: "من یک " + splitStr[randomNumber2] + " پیدا کردم",
					Status:      "found",
					Alert:       true,
					Chat:        true,
				},
				Addresses: []DTO2.CreateAddressDTO{
					{
						Province:      "تهران",
						City:          "تهران",
						AddressDetail: data.Districts[randomNumber].Name,
						Latitude:      data.Districts[randomNumber].Centroid.Latitude,
						Longitude:     data.Districts[randomNumber].Centroid.Longitude,
					},
				},
				ImgUrls: []string{},
				Tags: []string{
					tagNames[randomNumber4],
				},
			}
			posterRepository := Repository.NewPosterRepository(pgsql.GetDB())
			model, err := posterRepository.CreatePoster(uint64(userID), request.Poster, request.Addresses, nil,
				request.Tags, "normal")
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println(model)
		}
	}
	return nil
}

type CreateMockDataRequest struct {
	Count    int      `json:"count" binding:"required"`
	UserId   int      `json:"user_id" binding:"required"`
	TagNames []string `json:"tag_names" binding:"required"`
}

func CreateMockDataResponse(c *gin.Context) {
	var request CreateMockDataRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := MockPoster(request.Count, request.UserId, request.TagNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mock data created!"})
}
