package server

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/handler"

	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Структура Сервера
type Server struct {
	httpServer *http.Server
}

// Метод запускающий сервер
func (s *Server) Run(cfg *configs.Config, handler *handler.Handler) chan error {
	serverErr := make(chan error)
	s.httpServer = &http.Server{
		Addr:      cfg.SrvAddr,
		Handler:   handler.Router,
		TLSConfig: cfg.TlsConf,
	}
	if cfg.EnableHTTPS {
		go func() {
			logrus.Println("HTTPS Run")
			if err := s.httpServer.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				serverErr <- err
			}
		}()
		return serverErr
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	return serverErr
}

// Метод прекращающий работу сервера
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
