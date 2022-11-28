package storage

import (
	"fmt"
	"strconv"
	"sync"
)

type shortURL struct {
	Id     int
	Result string
	Origin string
}

type UrlStorage struct {
	sync.Mutex
	store  map[int]shortURL
	nextId int
}

func NewStorage() *UrlStorage {
	store := &UrlStorage{
		store:  make(map[int]shortURL),
		nextId: 0,
	}

	return store
}

func (store *UrlStorage) GetShortUrl(id int) (sUrl shortURL, er error) {
	store.Lock()
	defer store.Unlock()

	s, ok := store.store[id]
	if ok {
		return s, nil
	} else {
		return shortURL{}, fmt.Errorf("url with id=%d not found", id)
	}
}

func (store *UrlStorage) AddUrl(input string) string {
	store.Lock()
	defer store.Unlock()

	url := shortURL{
		Id:     store.nextId,
		Result: strconv.Itoa(store.nextId),
		Origin: input}
	// TODO написать функцию которая будет сокращать ссылку

	store.store[store.nextId] = url
	store.nextId++
	return url.Result
}
