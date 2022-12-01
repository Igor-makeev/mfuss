package handler

import (
	"fmt"
	"io"
	"mfuss/internal/repositories"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	repository *repositories.Repositories
}

func NewHandler(repository *repositories.Repositories) *Handler {

	return &Handler{repository: repository}
}

func (h *Handler) RootHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {

		switch req.Method {
		case http.MethodPost:
			h.PostHandler(w, req)

		default:
			http.Error(w, fmt.Sprintf("expect method POST at /, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}

	} else {
		switch req.Method {
		case http.MethodGet:

			pathParts := strings.Split(req.URL.Path, "/")

			if len(pathParts) < 2 {
				http.Error(w, "expect /<id> in  GetURLHandler handler", http.StatusBadRequest)
				return
			}
			id, err := strconv.Atoi(pathParts[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			h.GetURLHandler(w, req, id)
		default:
			http.Error(w, fmt.Sprintf("expect method GET at /{id}, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(string(b)); err != nil {

		http.Error(w, fmt.Sprintf("invalid URL: %v", string(b)), http.StatusInternalServerError)

		return
	}

	short := "http://" + r.Host + r.URL.Path + h.repository.SaveURL(string(b))

	if _, err := url.ParseRequestURI(short); err != nil {
		http.Error(w, fmt.Sprintf("output data: %v is invalid URL", short), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(short))

}

func (h *Handler) GetURLHandler(w http.ResponseWriter, req *http.Request, id int) {

	sURL, err := h.repository.GetShortURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", sURL.Origin)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
