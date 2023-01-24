package mock

import (
	"fmt"
	"mfuss/configs"
	"mfuss/internal/entity"
)

type StorageMock struct {
	store map[string]entity.ShortURL
	ID    string
	cfg   configs.Config
}

func NewStorageMock(cfg *configs.Config) *StorageMock {
	return &StorageMock{store: make(map[string]entity.ShortURL), ID: "0", cfg: *cfg}
}

func (store *StorageMock) SaveURL(input, userid string) (string, error) {
	url := entity.ShortURL{
		ID:     store.ID,
		Origin: input}

	store.store[store.ID] = url
	url.ResultURL = store.cfg.BaseURL + "/" + url.ID
	return url.ResultURL, nil
}
func (store *StorageMock) GetAllURLS(userID string) []entity.ShortURL {

	var urls []entity.ShortURL
	for _, v := range store.store {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}
	return urls
}

func (store *StorageMock) GetShortURL(id, idstring string) (sURL entity.ShortURL, er error) {
	s, ok := store.store[id]
	if ok {
		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}
func (store *StorageMock) Close() error {
	return nil
}
func (store *StorageMock) MultipleShort(input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	return nil, nil
}

func (store *StorageMock) MarkAsDeleted(arr []string, id string) {

}
func (store *StorageMock) Ping() error {
	return nil
}

type PersistentStorageMock struct {
}

func NewPersistentStorageMock() *PersistentStorageMock {
	return &PersistentStorageMock{}
}

func (psm *PersistentStorageMock) SaveData(ms map[string]entity.ShortURL) error {

	return nil

}

func (psm *PersistentStorageMock) LoadData(ms map[string]entity.ShortURL) error {
	return nil
}

func (psm *PersistentStorageMock) Close() error {
	return nil
}
