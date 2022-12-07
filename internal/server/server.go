package server

import (
	"context"
	"net/http"
)

type URLserver struct {
	server *http.Server
}

func NewURLServer(h http.Handler) *URLserver {
	server := URLserver{
		server: &http.Server{
			Addr:    ":8080",
			Handler: h,
		},
	}

	return &server
}

func (s *URLserver) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *URLserver) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
