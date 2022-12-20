package mock

import (
	"fmt"
	"mfuss/internal/entity"
)

type StorageMock struct {
	store map[string]entity.ShortURL
	ID    string
}

func NewStorageMock() *StorageMock {
	return &StorageMock{store: make(map[string]entity.ShortURL), ID: "0"}
}

func (store *StorageMock) SaveURL(input string) (string, error) {
	url := entity.ShortURL{
		ID:     store.ID,
		Origin: input}

	store.store[store.ID] = url

	return url.ID, nil
}

func (store *StorageMock) GetShortURL(id string) (sURL entity.ShortURL, er error) {
	s, ok := store.store[id]
	if ok {
		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}
func (store *StorageMock) Close() error {
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
