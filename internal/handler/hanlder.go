package handler

import (
	"mfuss/internal/repositories"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *repositories.Repositories
}

func NewHandler(repository *repositories.Repositories) *Handler {

	return &Handler{repository: repository}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	root := router.Group("/")
	{
		root.POST("/", h.PostHandler)
		root.GET("/:id", h.GetURLHandler)
	}

	return router
}
