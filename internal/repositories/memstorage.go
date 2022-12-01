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
	store := &MemoryStorage{
		store:  make(map[int]entity.ShortURL),
		nextID: 0,
	}

	return store
}

func (store *MemoryStorage) GetShortURL(id int) (sURL entity.ShortURL, er error) {
	store.Lock()
	defer store.Unlock()

	s, ok := store.store[id]
	if ok {
		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%d not found", id)

}

func (store *MemoryStorage) SaveURL(input string) string {
	store.Lock()
	defer store.Unlock()

	url := entity.ShortURL{
		ID:     store.nextID,
		Result: strconv.Itoa(store.nextID),
		Origin: input}

	store.store[store.nextID] = url
	store.nextID++
	return url.Result
}
