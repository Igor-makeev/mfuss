package repositories

import (
	"fmt"
	"math/rand"
	"mfuss/internal/entity"
	"sync"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type MemoryStorage struct {
	sync.Mutex
	Store map[string]entity.ShortURL
	PersistentStorage
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Store: make(map[string]entity.ShortURL),
	}

}

func (ms *MemoryStorage) GetShortURL(id string) (sURL entity.ShortURL, er error) {
	ms.Lock()
	defer ms.Unlock()

	s, ok := ms.Store[id]
	if ok {
		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}

func (ms *MemoryStorage) SaveURL(input string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	url := entity.ShortURL{
		ID:     genetareID(),
		Origin: input}

	ms.Store[url.ID] = url

	return url.ID, nil
}

func genetareID() string {
	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	res := string(buf)
	return res
}
