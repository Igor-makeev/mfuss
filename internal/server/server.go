package server

import (
	"context"
	"net/http"
)

type URLserver struct {
	server *http.Server
}

func (s *URLserver) Run(h http.Handler) error {
	s.server = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	return s.server.ListenAndServe()
}

func (s *URLserver) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
