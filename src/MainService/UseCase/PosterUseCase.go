package UseCase

import (
	"encoding/json"
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	DTO2 "github.com/403-access-denied/main-backend/src/MainService/DTO"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/403-access-denied/main-backend/src/MainService/View"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type getPostersRequest struct {
	PageID       int     `form:"page_id" binding:"required,min=1"`
	PageSize     int     `form:"page_size" binding:"required,min=5"`
	Sort         string  `form:"sort,omitempty" binding:"omitempty,oneof=asc desc"`
	SortBy       string  `form:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at id"`
	Status       string  `form:"status,omitempty" binding:"oneof=lost found both"`
	SearchPhrase string  `form:"search_phrase,omitempty"`
	TimeStart    int64   `form:"time_start,omitempty"`
	TimeEnd      int64   `form:"time_end,omitempty"`
	OnlyRewards  bool    `form:"only_rewards,omitempty"`
	Lat          float64 `form:"lat,omitempty"`
	Lon          float64 `form:"lon,omitempty"`
	TagIds       []int   `form:"tag_ids,omitempty" swaggertype:"array,int"`
	State        string  `form:"state,omitempty" binding:"omitempty,oneof=all pending accepted rejected"`
}

func GetPostersResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())

	var request getPostersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset := (request.PageID - 1) * request.PageSize
	request.Sort = c.DefaultQuery("sort", "asc")
	request.SortBy = c.DefaultQuery("sort_by", "created_at")
	request.Status = c.DefaultQuery("status", "both")
	request.SearchPhrase = c.DefaultQuery("search_phrase", "")
	//todo add other fields

	filterObject := DTO2.FilterObject{
		Status:       request.Status,
		SearchPhrase: request.SearchPhrase,
		TimeStart:    request.TimeStart,
		TimeEnd:      request.TimeEnd,
		OnlyRewards:  request.OnlyRewards,
		Lat:          request.Lat,
		Lon:          request.Lon,
		TagIds:       request.TagIds,
		State:        request.State,
	}

	posters, err := posterRepository.GetAllPosters(request.PageSize, offset, request.Sort, request.SortBy, filterObject)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//DBConfiguration.CloseDB()
	View.GetPostersView(posters, c)
}

type GetPosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func GetPosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
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
	//DBConfiguration.CloseDB()
	View.GetPosterByIdView(poster, c)
}

type DeletePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func DeletePosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
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
	ImgUrls   []string `json:"img_urls" binding:"required"`
	Tags      []string `json:"tags" binding:"required"`
}

func CreatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	var request CreatePosterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Poster.UserID = uint(payload.UserID)
	poster, err := posterRepository.CreatePoster(request.Poster, request.Addresses, request.ImgUrls, request.Tags)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//DBConfiguration.CloseDB()
	View.CreatePosterView(poster, c)
}

type UpdatePosterRequest struct {
	Poster    DTO2.UpdatePosterDTO    `json:"poster"`
	Addresses []DTO2.UpdateAddressDTO `json:"addresses"`
}

type UpdatePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func UpdatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	payload := c.MustGet("authorization_payload").(*Token.Payload)

	var request UpdatePosterRequest
	var id UpdatePosterByIdRequest
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Poster.UserID = uint(payload.UserID)
	err := posterRepository.UpdatePoster(id.ID, request.Poster, request.Addresses)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poster updated successfully"})
}

type GetPhotoAiNSFWRequest struct {
	ImageUrl string `json:"image_url" binding:"required"`
}

func GetPhotoAiNSFWResponse(c *gin.Context) {
	var request GetPhotoAiNSFWRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	boolean := false
	_ = boolean
	PhotoUrl := request.ImageUrl
	fmt.Println(PhotoUrl)
	channel := make(chan bool)
	go func() {

		url := fmt.Sprintf("https://api.apilayer.com/nudity_detection/url?url=%s", PhotoUrl)
		fmt.Println(url)
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("apikey", "z232GHVwPAec88LzsqdBjUhL5BZDgvGp")

		if err != nil {
			fmt.Println(err)
		}
		res, err := client.Do(req)
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, err := ioutil.ReadAll(res.Body)
		splitStr := strings.Split(string(body), ",")
		splitStr2 := strings.Split(splitStr[0], ": ")
		a, err := strconv.Atoi(splitStr2[1])
		fmt.Println(a)
		if a > 1 {
			boolean = true
		}
		channel <- true
	}()
	fmt.Println("continue ...")
	<-channel
	if boolean {
		c.JSON(http.StatusOK, gin.H{"message": "nsfw"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "sfw"})
	}
}

type GetTextNSFWRequest struct {
	Text string `form:"text" binding:"required"`
}

func GetTextNSFWResponse(c *gin.Context) {
	var request GetTextNSFWRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, err := ioutil.ReadFile("Utils/data.txt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	boolean := false
	_ = boolean
	newString := strings.ReplaceAll(string(file), "\n", "")
	newString2 := strings.ReplaceAll(newString, "\"", "")
	newString3 := strings.ReplaceAll(newString2, "\r", "")
	splitStr := strings.Split(newString3, ",")
	text := request.Text
	for i, _ := range splitStr {
		if strings.Contains(text, splitStr[i]) {
			boolean = true
		}
	}
	if boolean {
		c.JSON(http.StatusOK, gin.H{"message": "nsfw"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "sfw"})
	}
}

type CreatePosterReportRequest struct {
	PosterId    uint   `form:"poster_id" binding:"required"`
	IssuerId    uint   `form:"issuer_id" binding:"required"`
	ReportType  string `form:"report_type" binding:"required,oneof=spam inappropriate other"` //TODO: add more report types
	Description string `form:"description"`
}

func CreatePosterReportResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
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
	//DBConfiguration.CloseDB()
	c.JSON(http.StatusOK, gin.H{"message": "Report created"})
}

type getPosterReportsRequest struct {
	PageID   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=20"`
	Status   string `form:"status,omitempty" binding:"oneof=open resolved both"`
}

func GetPosterReportsResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())

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

	//DBConfiguration.CloseDB()
	View.GetPosterReportsView(posterReports, c)
}

type GetPosterReportByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func GetPosterReportByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())

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
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())

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
	//DBConfiguration.CloseDB()

	c.JSON(http.StatusOK, gin.H{"message": "Report resolved"})
}

func UploadPosterImageResponse(c *gin.Context) {
	formHeader, err := c.FormFile("poster_image")
	fileName := formHeader.Filename
	extension := path.Ext(fileName)

	currentTime := time.Now().Format("20060102_150405")
	randomString := strconv.FormatInt(rand.Int63(), 16)
	newName := currentTime + "_" + randomString + extension
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := formHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	uploadUrl, err := Utils.UploadInArvanCloud(file, newName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": uploadUrl})
}

type UpdatePosterStateRequest struct {
	ID    uint   `form:"id" binding:"required,min=1"`
	State string `form:"state" binding:"required,oneof=accepted rejected pending"`
}

func UpdatePosterStateResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())

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
	//DBConfiguration.CloseDB()

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
					UserID:      uint(userID),
					State:       "pending",
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
			posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
			model, err := posterRepository.CreatePoster(request.Poster, request.Addresses, nil, request.Tags)
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
					UserID:      uint(userID),
					State:       "pending",
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
			posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
			model, err := posterRepository.CreatePoster(request.Poster, request.Addresses, nil, request.Tags)
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
