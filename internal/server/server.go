package server

import (
	"net/http"
)

func NewURLServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

}
