package repositories

import (
	"fmt"
	"mfuss/internal/entity"
	"strconv"
	"sync"
)

type MemoryStorage struct {
	sync.Mutex
	store  map[int]entity.ShortURL
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		store:  make(map[int]entity.ShortURL),
		nextID: 0,
	}

}

func (ms *MemoryStorage) GetShortURL(id int) (sURL entity.ShortURL, er error) {
	ms.Lock()
	defer ms.Unlock()

	s, ok := ms.store[id]
	if ok {
		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%d not found", id)

}

func (ms *MemoryStorage) SaveURL(input string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	url := entity.ShortURL{
		ID:     ms.nextID,
		Result: strconv.Itoa(ms.nextID),
		Origin: input}

	ms.store[ms.nextID] = url
	ms.nextID++
	return url.Result, nil
}
