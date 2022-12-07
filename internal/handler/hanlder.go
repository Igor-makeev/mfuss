package handler

import (
	"mfuss/internal/repositories"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *repositories.Repositories
	Router     *gin.Engine
}

func NewHandler(repository *repositories.Repositories) *Handler {
	handler := &Handler{
		Router:     gin.New(),
		repository: repository,
	}
	root := handler.Router.Group("/")
	{
		root.POST("/", handler.PostHandler)
		root.GET("/:id", handler.GetURLHandler)
	}

	return handler
}
