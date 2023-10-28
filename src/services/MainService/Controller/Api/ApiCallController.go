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
