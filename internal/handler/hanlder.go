package handler

import (
	"mfuss/internal/repositories"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo   *repositories.Repository
	Router *gin.Engine
}

func NewHandler(rep *repositories.Repository) *Handler {
	handler := &Handler{
		Router: gin.New(),
		Repo:   rep,
	}

	root := handler.Router.Group("/")
	{
		root.POST("/", handler.PostHandler)
		root.GET("/:id", handler.GetURLHandler)

	}

	api := handler.Router.Group("/api")
	{
		api.POST("/shorten", handler.PostJSONHandler)
	}

	return handler
}
