package service

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/repositories"
)

type Service struct {
	URLStorager
	Queue *Queue
	Cfg   *configs.Config
}

func NewService(repos *repositories.Repository) *Service {
	return &Service{
		URLStorager: NewURLStorageService(repos.URLStorager),
		Queue:       NewQueue(repos),
		Cfg:         repos.Config,
	}
}

func (service *Service) Close(ctx context.Context) error {
	if err := service.URLStorager.Close(ctx); err != nil {
		return err
	}
	service.Queue.Close()
	return nil
}
