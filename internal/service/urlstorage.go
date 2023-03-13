package service

import (
	"context"
	"mfuss/internal/entity"
	"mfuss/internal/repositories"
)

// интерфейс сервиса хранилища ссылок
type URLStorager interface {
	SaveURL(ctx context.Context, input, userID string) (string, error)
	GetAllURLs(ctx context.Context, userID string) []entity.ShortURL
	GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error)
	MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error)
	MarkAsDeleted(ctx context.Context, arr []string) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

// структура сервиса хранилища ссылок
type URLStorageService struct {
	repo repositories.URLStorager
}

// конструктор
func NewURLStorageService(repo repositories.URLStorager) *URLStorageService {
	return &URLStorageService{repo: repo}
}

// сохранить
func (uss *URLStorageService) SaveURL(ctx context.Context, input, userID string) (string, error) {

	return uss.repo.SaveURL(ctx, input, userID)
}

// получить все ссылки
func (uss *URLStorageService) GetAllURLs(ctx context.Context, userID string) []entity.ShortURL {
	return uss.repo.GetAllURLs(ctx, userID)
}

// получить ссылку
func (uss *URLStorageService) GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error) {
	return uss.repo.GetShortURL(ctx, id, userID)
}

// сокращение пачки ссылок
func (uss *URLStorageService) MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	return uss.repo.MultipleShort(ctx, input, userID)
}

// пометить ссылки как удаленные
func (uss *URLStorageService) MarkAsDeleted(ctx context.Context, arr []string) error {
	return uss.repo.MarkAsDeleted(ctx, arr)
}

// проверка соединения с бд
func (uss *URLStorageService) Ping(ctx context.Context) error {
	return uss.repo.Ping(ctx)
}

// закрыть
func (uss *URLStorageService) Close(ctx context.Context) error {
	return uss.repo.Close(ctx)
}
