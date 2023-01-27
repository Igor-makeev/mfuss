package service

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/repositories"
)

type URLStorager interface {
	SaveURL(input, userID string, ctx context.Context) (string, error)
	GetAllURLs(userID string, ctx context.Context) []entity.ShortURL
	GetShortURL(id, userID string, ctx context.Context) (sURL entity.ShortURL, er error)
	MultipleShort(input []entity.URLBatchInput, userID string, ctx context.Context) ([]entity.URLBatchResponse, error)
	MarkAsDeleted(arr []string, ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

type Service struct {
	URLStorager
	Queue *Queue
	Cfg   *configs.Config
}

func NewService(repos *repositories.Repository) *Service {
	return &Service{
		URLStorager: NewURLStorageService(repos.URLStorager),
		Queue:       NewQueue(),
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
