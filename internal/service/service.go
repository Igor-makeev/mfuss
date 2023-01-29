package service

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/repositories"
)

type URLStorager interface {
	SaveURL(ctx context.Context, input, userID string) (string, error)
	GetAllURLs(ctx context.Context, userID string) []entity.ShortURL
	GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error)
	MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error)
	MarkAsDeleted(ctx context.Context, arr []string) error
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
