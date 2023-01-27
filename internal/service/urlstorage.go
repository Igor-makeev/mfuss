package service

import (
	"context"
	"mfuss/internal/entity"
	"mfuss/internal/repositories"
)

type URLStorageService struct {
	repo repositories.URLStorager
}

func NewURLStorageService(repo repositories.URLStorager) *URLStorageService {
	return &URLStorageService{repo: repo}
}

func (uss *URLStorageService) SaveURL(input, userID string, ctx context.Context) (string, error) {

	return uss.repo.SaveURL(input, userID, ctx)
}
func (uss *URLStorageService) GetAllURLs(userID string, ctx context.Context) []entity.ShortURL {
	return uss.repo.GetAllURLs(userID, ctx)
}
func (uss *URLStorageService) GetShortURL(id, userID string, ctx context.Context) (sURL entity.ShortURL, er error) {
	return uss.repo.GetShortURL(id, userID, ctx)
}
func (uss *URLStorageService) MultipleShort(input []entity.URLBatchInput, userID string, ctx context.Context) ([]entity.URLBatchResponse, error) {
	return uss.repo.MultipleShort(input, userID, ctx)
}
func (uss *URLStorageService) MarkAsDeleted(arr []string, ctx context.Context) error {
	return uss.repo.MarkAsDeleted(arr, ctx)
}
func (uss *URLStorageService) Ping(ctx context.Context) error {
	return uss.repo.Ping(ctx)
}
func (uss *URLStorageService) Close(ctx context.Context) error {
	return uss.repo.Close(ctx)
}
