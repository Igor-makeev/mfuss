package server

import (
	"mfuss/internal/handler"
	"net/http"
)

func NewURLServer(h *handler.Handler) *http.Server {
	return &http.Server{
		Addr:    h.Service.Cfg.SrvAddr,
		Handler: h.Router,
	}

}
