package handler

import (
	"compress/gzip"
	"mfuss/internal/service"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Стркуктура Обработчика
type Handler struct {
	Service *service.Service
	Router  *gin.Engine
}

// Конструктор обработчика
func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		Router:  gin.New(),
		Service: service,
	}

	pprof.Register(handler.Router)

	root := handler.Router.Group("/", handler.userCheck)
	{

		root.POST("/", GzipUnpack(), handler.PostHandler)
		root.GET("/:id", GzipCompress(gzip.DefaultCompression), handler.GetURLHandler)
		root.GET("/ping", handler.GetPingHandler)

		api := root.Group("api")
		{
			api.POST("/shorten/batch", handler.MultipleShortHandler)
			api.POST("/shorten", GzipUnpack(), handler.PostJSONHandler)
			api.GET("/user/urls", handler.GetUserURLs)
			api.DELETE("/user/urls", handler.checkURLSID, handler.DeleteUrls)
		}
	}

	return handler
}
