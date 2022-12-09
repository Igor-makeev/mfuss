package handler

import (
	"mfuss/internal/repositories"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage repositories.URLStorage
	Router  *gin.Engine
}

func NewHandler(ms repositories.URLStorage) *Handler {
	handler := &Handler{
		Router:  gin.New(),
		storage: ms,
	}
	root := handler.Router.Group("/")
	{
		root.POST("/", handler.PostHandler)
		root.GET("/:id", handler.GetURLHandler)
	}

	return handler
}
