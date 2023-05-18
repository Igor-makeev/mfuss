package mock

import (
	"context"
	"fmt"
	"mfuss/configs"
	"mfuss/internal/entity"
)

// Мок
type StorageMock struct {
	store map[string]*entity.ShortURL
	ID    string
	cfg   configs.Config
}

// Мок
func NewStorageMock(cfg *configs.Config) *StorageMock {
	return &StorageMock{store: make(map[string]*entity.ShortURL), ID: "0", cfg: *cfg}
}

// Мок
func (store *StorageMock) SaveURL(ctx context.Context, input, userid string) (string, error) {

	url := entity.ShortURL{
		ID:     store.ID,
		Origin: input}

	store.store[store.ID] = &url
	url.ResultURL = store.cfg.BaseURL + "/" + url.ID
	return url.ResultURL, nil
}

// Мок
func (store *StorageMock) GetAllURLs(ctx context.Context, userID string) []entity.ShortURL {

	var urls []entity.ShortURL
	for _, v := range store.store {
		if v.UserID == userID {
			urls = append(urls, *v)
		}
	}
	return urls
}

// Мок
func (store *StorageMock) GetShortURL(ctx context.Context, id, idstring string) (sURL entity.ShortURL, er error) {
	s, ok := store.store[id]
	if ok {
		return *s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}

// Мок
func (store *StorageMock) Close(ctx context.Context) error {
	return nil
}

// Мок
func (store *StorageMock) MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	return nil, nil
}

// Мок
func (store *StorageMock) MarkAsDeleted(ctx context.Context, arr []string) error {
	return nil
}

// Мок
func (store *StorageMock) Ping(ctx context.Context) error {
	return nil
}

// Мок
func (store *StorageMock) GetStats(ctx context.Context) (entity.Stats, error) {
	return entity.Stats{}, nil
}

// Мок
type PersistentStorageMock struct {
}

// Мок
func NewPersistentStorageMock() *PersistentStorageMock {
	return &PersistentStorageMock{}
}

// Мок
func (psm *PersistentStorageMock) SaveData(ms map[string]entity.ShortURL) error {

	return nil

}

// Мок
func (psm *PersistentStorageMock) LoadData(ms map[string]entity.ShortURL) error {
	return nil
}

// Мок
func (psm *PersistentStorageMock) Close(ctx context.Context) error {
	return nil
}
