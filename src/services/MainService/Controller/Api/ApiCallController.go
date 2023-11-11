package Api

import (
	"github.com/200-status-ok/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// ImageUpload godoc
// @Summary Upload image
// @Description Upload image
// @Tags ApiCall
// @Accept  multipart/form-data
// @Param files formData file true "Multiple files"
// @Produce  json
// @Success 200
// @Router /api-call/image-upload [post]
func ImageUpload(c *gin.Context) {
	UseCase.ImageUploadResponse(c)
}

// GeneratePosterInfo godoc
// @Summary Generate poster info
// @Description Generates info for a poster
// @Tags ApiCall
// @Accept  json
// @Produce  json
// @Param image_url query string true "Image Url"
// @Success 200 {object} View.GeneratedPosterInfoView
// @Router /api-call/generate-poster-Info [get]
func GeneratePosterInfo(c *gin.Context) {
	UseCase.GeneratePosterInfoResponse(c)
}
