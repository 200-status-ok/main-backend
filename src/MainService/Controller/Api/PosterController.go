package Api

import (
	"github.com/403-access-denied/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// GetPosters godoc
// @Summary Get a list of all posters
// @Description Retrieves a list of all posters, sorted and paginated according to the given parameters
// @Tags posters
// @Accept  json
// @Produce  json
// @Param page_id query int true "Page ID" minimum(1) default(1)
// @Param page_size query int true "Page size" minimum(1) default(10)
// @Param sort query string false "Sort direction" enum(asc, desc) default(asc)
// @Param sort_by query string false "Sort by" enum(id, updated_at, created_at) default(created_at)
// @Param search_phrase query string false "Search phrase"
// @Param status query string false "Status" enum(lost, found, both) default(both)
// @Param time_start query int false "Time start"
// @Param time_end query int false "Time end"
// @Param only_rewards query bool false "Only rewards"
// @Param lat query float64 false "Latitude"
// @Param lon query float64 false "Longitude"
// @Param tag_ids query []int false "TagIds" collectionFormat(multi) example(1,2,3)
// @Param state query string false "State" enum(all, accepted, rejected, pending) default(all)
// @Success 200 {array} View.PosterView
// @Router /posters [get]
func GetPosters(c *gin.Context) {
	UseCase.GetPostersResponse(c)
}

// GetPoster godoc
// @Summary Get a poster by ID
// @Description Retrieves a poster by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200 {object} View.PosterView
// @Router /posters/{id} [get]
func GetPoster(c *gin.Context) {
	UseCase.GetPosterByIdResponse(c)
}

// CreatePoster godoc
// @Summary Create a poster
// @Description Creates a poster
// @Tags posters
// @Accept  json
// @Produce  json
// @Param poster body UseCase.CreatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters/authorize [post]
func CreatePoster(c *gin.Context) {
	UseCase.CreatePosterResponse(c)
}

// UpdatePoster godoc
// @Summary Update a poster by ID
// @Description Updates a poster by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Param poster body UseCase.UpdatePosterRequest true "Poster"
// @Success 200 {object} View.PosterView
// @Router /posters/authorize/{id} [patch]
func UpdatePoster(c *gin.Context) {
	UseCase.UpdatePosterResponse(c)
}

// DeletePoster godoc
// @Summary Delete a poster by ID
// @Description Deletes a poster by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id path int true "Poster ID"
// @Success 200
// @Router /posters/authorize/{id} [delete]
func DeletePoster(c *gin.Context) {
	UseCase.DeletePosterByIdResponse(c)
}

// GetPhotoNSFWAi godoc
// @Summary Get photo nsfw AI
// @Description Get photo nsfw AI
// @Tags ApiCall
// @Accept  json
// @Produce  json
// @Param image_url body UseCase.GetPhotoAiNSFWRequest true "Image Url"
// @Success 200
// @Router /api-call/predict [post]
func GetPhotoNSFWAi(c *gin.Context) {
	UseCase.GetPhotoAiNSFWResponse(c)
}

// GetTextNSFW godoc
// @Summary Get text nsfw
// @Description Get text nsfw
// @Tags ApiCall
// @Accept  json
// @Produce  json
// @Param text query string true "Text"
// @Success 200
// @Router /api-call/predict-txt [get]
func GetTextNSFW(c *gin.Context) {
	UseCase.GetTextNSFWResponse(c)
}

// CreatePosterReport godoc
// @Summary Report a poster
// @Description Reports a poster
// @Tags reports
// @Accept  json
// @Produce  json
// @Param poster_id query int true "Poster ID"
// @Param issuer_id query int true "Issuer ID"
// @Param report_type query string true "Report Type" enum(spam, inappropriate, other) default(other)
// @Param description query string false "Description"
// @Success 200
// @Router /reports/report-poster [post]
func CreatePosterReport(c *gin.Context) {
	UseCase.CreatePosterReportResponse(c)
}

// GetPosterReports godoc
// @Summary Get a list of all poster reports
// @Description Retrieves a list of all poster reports, sorted and paginated according to the given parameters
// @Tags reports
// @Accept  json
// @Produce  json
// @Param page_id query int true "Page ID" minimum(1) default(1)
// @Param page_size query int true "Page size" minimum(1) default(10)
// @Param status query string true "Status" enum(open, resolved, both) default(both)
// @Success 200 {array} View.PosterReportView
// @Router /reports [get]
func GetPosterReports(c *gin.Context) {
	UseCase.GetPosterReportsResponse(c)
}

// GetPosterReport godoc
// @Summary Get a poster report by ID
// @Description Retrieves a poster report by ID
// @Tags reports
// @Accept  json
// @Produce  json
// @Param id path int true "Report ID"
// @Success 200 {object} View.PosterReportView
// @Router /reports/{id} [get]
func GetPosterReport(c *gin.Context) {
	UseCase.GetPosterReportByIdResponse(c)
}

// UpdatePosterReport godoc
// @Summary Update a poster report by ID
// @Description Updates a poster report by ID
// @Tags reports
// @Accept  json
// @Produce  json
// @Param id path int true "Report ID"
// @Param report body UseCase.UpdatePosterReportRequest true "Poster Report"
// @Success 200 {object} View.PosterView
// @Router /reports/{id} [patch]
func UpdatePosterReport(c *gin.Context) {
	UseCase.UpdatePosterReportResponse(c)
}

// UploadPosterImage godoc
// @Summary Upload poster image
// @Description Upload poster image
// @Tags posters
// @Accept  multipart/form-data
// @Param files formData file true "Multiple files"
// @Produce  json
// @Success 200
// @Router /posters/upload-image [post]
func UploadPosterImage(c *gin.Context) {
	UseCase.UploadPosterImageResponse(c)
}

// UpdatePosterState godoc
// @Summary Update a poster state by ID
// @Description Updates a poster report by ID
// @Tags posters
// @Accept  json
// @Produce  json
// @Param id query int true "ID"
// @Param state query string false "State" enum(pending, accepted, rejected) default(accepted)
// @Success 200
// @Router /posters/update-state [patch]
func UpdatePosterState(c *gin.Context) {
	UseCase.UpdatePosterStateResponse(c)
}
