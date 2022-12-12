package handler

import (
	"mfuss/configs"
	"mfuss/internal/repositories"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage repositories.URLStorage
	Router  *gin.Engine
	Cfg     configs.Config
}

func NewHandler(ms repositories.URLStorage, cfg configs.Config) *Handler {
	handler := &Handler{
		Router:  gin.New(),
		storage: ms,
		Cfg:     cfg,
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
