package Controller

import (
	"github.com/403-access-denied/main-backend/src/Controller/Api"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func (s *Server) MainController() {
	v1 := s.Router.Group("/api/v1")
	{
		poster := v1.Group("/posters")
		{
			poster.GET("/", Api.GetPosters)
			poster.GET("/:id", Api.GetPoster)
			poster.POST("/", Api.CreatePoster)
			poster.PATCH("/:id", Api.UpdatePoster)
			poster.DELETE("/:id", Api.DeletePoster)
		}
	}

}
