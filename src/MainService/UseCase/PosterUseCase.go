package UseCase

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	DTO2 "github.com/403-access-denied/main-backend/src/MainService/DTO"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/403-access-denied/main-backend/src/MainService/View"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type getPostersRequest struct {
	PageID       int     `form:"page_id" binding:"required,min=1"`
	PageSize     int     `form:"page_size" binding:"required,min=5,max=20"`
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
	}

	posters, err := posterRepository.GetAllPosters(request.PageSize, offset, request.Sort, request.SortBy, filterObject)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DBConfiguration.CloseDB()
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
	DBConfiguration.CloseDB()
	View.GetPosterByIdView(poster, c)
}

type DeletePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func DeletePosterByIdResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request DeletePosterByIdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := posterRepository.DeletePosterById(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	c.JSON(http.StatusOK, gin.H{"message": "Poster deleted"})
}

type CreatePosterRequest struct {
	Poster    DTO2.PosterDTO
	Addresses []DTO2.AddressDTO
	ImgUrls   []string `json:"img_urls" binding:"required"`
	Tags      []int    `json:"tags" binding:"required"`
}

func CreatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	var request CreatePosterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poster, err := posterRepository.CreatePoster(request.Poster, request.Addresses, request.ImgUrls, request.Tags)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.CreatePosterView(poster, c)
}

type UpdatePosterRequest struct {
	Poster     DTO2.PosterDTO
	Addresses  []DTO2.AddressDTO
	ImgUrls    []string `json:"img_urls" binding:"required"`
	Categories []int    `json:"categories" binding:"required"`
}

type UpdatePosterByIdRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func UpdatePosterResponse(c *gin.Context) {
	posterRepository := Repository.NewPosterRepository(DBConfiguration.GetDB())
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
	poster, err := posterRepository.UpdatePoster(id.ID, request.Poster, request.Addresses, request.ImgUrls,
		request.Categories)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DBConfiguration.CloseDB()
	View.UpdatePosterView(poster, c)
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
	DBConfiguration.CloseDB()
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

	DBConfiguration.CloseDB()
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

	DBConfiguration.CloseDB()
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
	DBConfiguration.CloseDB()

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
