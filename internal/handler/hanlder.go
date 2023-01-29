package handler

import (
	"compress/gzip"
	"mfuss/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
	Router  *gin.Engine
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		Router:  gin.New(),
		Service: service,
	}
	handler.Router.Use(UserCheck())
	root := handler.Router.Group("/")
	{

		root.POST("/", GzipUnpack(), handler.PostHandler)
		root.GET("/:id", GzipCompress(gzip.DefaultCompression), handler.GetURLHandler)
		root.GET("/ping", handler.GetPingHandler)

		api := handler.Router.Group("api")
		{
			api.POST("/shorten/batch", handler.MultipleShortHandler)
			api.POST("/shorten", GzipUnpack(), handler.PostJSONHandler)
			api.GET("/user/urls", handler.GetUserURLs)
			api.DELETE("/user/urls", URLSIDCheck(), handler.DeleteUrls)
		}
	}

	return handler
}
