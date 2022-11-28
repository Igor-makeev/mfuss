package server

import (
	"context"
	"net/http"
)

type UrlServer struct {
	server *http.Server
}

func (s *UrlServer) Run(h http.Handler) error {
	s.server = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	return s.server.ListenAndServe()
}

func (s *UrlServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
