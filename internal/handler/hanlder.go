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

// func (h *Handler) RootHandler(w http.ResponseWriter, req *http.Request) {
// 	if req.URL.Path == "/" {

// 		switch req.Method {
// 		case http.MethodPost:
// 			h.PostHandler(w, req)

// 		default:
// 			http.Error(w, fmt.Sprintf("expect method POST at /, got %v", req.Method), http.StatusMethodNotAllowed)
// 			return
// 		}

// 	} else {
// 		switch req.Method {
// 		case http.MethodGet:

// 			pathParts := strings.Split(req.URL.Path, "/")

// 			if len(pathParts) < 2 {
// 				http.Error(w, "expect /<id> in  GetURLHandler handler", http.StatusBadRequest)
// 				return
// 			}
// 			id, err := strconv.Atoi(pathParts[1])
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 				return
// 			}
// 			h.GetURLHandler(w, req, id)
// 		default:
// 			http.Error(w, fmt.Sprintf("expect method GET at /{id}, got %v", req.Method), http.StatusMethodNotAllowed)
// 			return
// 		}
// 	}
// }
