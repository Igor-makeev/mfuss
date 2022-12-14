package server

import (
	"mfuss/internal/handler"
	"net/http"
)

func NewURLServer(h *handler.Handler) *http.Server {
	return &http.Server{
		Addr:    h.Repo.Config.SrvAddr,
		Handler: h.Router,
	}

}
