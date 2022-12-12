package server

import (
	"mfuss/configs"
	"net/http"
)

func NewURLServer(h http.Handler, cfg configs.Config) *http.Server {
	return &http.Server{
		Addr:    cfg.SrvAddr,
		Handler: h,
	}

}
