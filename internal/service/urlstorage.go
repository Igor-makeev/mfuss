package service

import (
	"context"
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

type URLStorageService struct {
	repo repositories.URLStorager
}

func NewURLStorageService(repo repositories.URLStorager) *URLStorageService {
	return &URLStorageService{repo: repo}
}

func (uss *URLStorageService) SaveURL(ctx context.Context, input, userID string) (string, error) {

	return uss.repo.SaveURL(ctx, input, userID)
}
func (uss *URLStorageService) GetAllURLs(ctx context.Context, userID string) []entity.ShortURL {
	return uss.repo.GetAllURLs(ctx, userID)
}
func (uss *URLStorageService) GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error) {
	return uss.repo.GetShortURL(ctx, id, userID)
}
func (uss *URLStorageService) MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	return uss.repo.MultipleShort(ctx, input, userID)
}
func (uss *URLStorageService) MarkAsDeleted(ctx context.Context, arr []string) error {
	return uss.repo.MarkAsDeleted(ctx, arr)
}
func (uss *URLStorageService) Ping(ctx context.Context) error {
	return uss.repo.Ping(ctx)
}
func (uss *URLStorageService) Close(ctx context.Context) error {
	return uss.repo.Close(ctx)
}
