package handler

import (
	"compress/gzip"
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

		root.POST("/", GzipUnpack(), handler.PostHandler)
		root.GET("/:id", GzipCompress(gzip.DefaultCompression), handler.GetURLHandler)

	}

	api := handler.Router.Group("/api")
	{
		api.POST("/shorten", GzipUnpack(), handler.PostJSONHandler)
	}

	return handler
}
