package Api

import (
	"github.com/403-access-denied/main-backend/src/MainService/UseCase"
	"github.com/gin-gonic/gin"
)

// SignupAdmin godoc
// @Summary signup admin
// @Description signup admin
// @Tags admins
// @Accept  json
// @Produce  json
// @Param admin body UseCase.SignupAdminRequest true "Signup Admin"
// @Success 200 {object} View.AdminView
// @Router /admin/signup [post]
func SignupAdmin(c *gin.Context) {
	UseCase.SignupAdminResponse(c)
}

// LoginAdmin godoc
// @Summary login admin
// @Description login admin
// @Tags admins
// @Accept  json
// @Produce  json
// @Param admin body UseCase.LoginAdminRequest true "Login Admin"
// @Success 200 {object} View.AdminLoginView
// @Router /admin/login [post]
func LoginAdmin(c *gin.Context) {
	UseCase.LoginAdminResponse(c)
}
