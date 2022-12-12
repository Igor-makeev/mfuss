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
