package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mfuss/internal/storage"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type MyHandler struct {
	store *storage.URLStorage
}

func NewHandler() *MyHandler {
	handler := &MyHandler{
		store: storage.NewStorage(),
	}
	return handler
}

func (h *MyHandler) RootHandler(w http.ResponseWriter, req *http.Request) {
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

func (h *MyHandler) PostHandler(w http.ResponseWriter, r *http.Request) {

	logrus.Printf("handling AddURL  %s\n", r.URL.Path)

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	short := r.Host + r.URL.Path + h.store.AddURL(string(b))

	if _, err := url.Parse(short); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(js)

}

func (h *MyHandler) GetURLHandler(w http.ResponseWriter, req *http.Request, id int) {
	logrus.Printf("handling GetURL at %s\n", req.URL.Path)

	sURL, err := h.store.GetShortURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", sURL.Origin)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(sURL.Result))

}
