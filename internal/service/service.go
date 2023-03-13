package service

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/repositories"
)

// структура сервиса
type Service struct {
	URLStorager
	Queue *Queue
	Cfg   *configs.Config
}

// констурктор сервиса
func NewService(repos *repositories.Repository) *Service {
	return &Service{
		URLStorager: NewURLStorageService(repos.URLStorager),
		Queue:       NewQueue(repos),
		Cfg:         repos.Config,
	}
}

// метод закрытия сервиса
func (service *Service) Close(ctx context.Context) error {
	if err := service.URLStorager.Close(ctx); err != nil {
		return err
	}
	service.Queue.Close()
	return nil
}
