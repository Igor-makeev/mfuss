package storage

import (
	"fmt"
	"strconv"
	"sync"
)

type shortURL struct {
	ID     int
	Result string
	Origin string
}

type URLStorage struct {
	sync.Mutex
	store  map[int]shortURL
	nextID int
}

const key = "http://localhost:8080/"

func NewStorage() *URLStorage {
	store := &URLStorage{
		store:  make(map[int]shortURL),
		nextID: 0,
	}

	return store
}

func (store *URLStorage) GetShortURL(id int) (sURL shortURL, er error) {
	store.Lock()
	defer store.Unlock()

	s, ok := store.store[id]
	if ok {
		return s, nil
	} else {
		return shortURL{}, fmt.Errorf("url with id=%d not found", id)
	}
}

func (store *URLStorage) AddURL(input string) string {
	store.Lock()
	defer store.Unlock()

	url := shortURL{
		ID:     store.nextID,
		Result: key + strconv.Itoa(store.nextID),
		Origin: input}

	store.store[store.nextID] = url
	store.nextID++
	return url.Result
}
