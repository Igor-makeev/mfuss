package repositories

import (
	"fmt"
	"math/rand"
	"mfuss/configs"
	"mfuss/internal/entity"
	"sync"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type PersistentStorage interface {
	SaveData(ms map[string]entity.ShortURL) error
	LoadData(ms map[string]entity.ShortURL) error
	Close() error
}

type MemoryStorage struct {
	sync.Mutex
	URLStore map[string]entity.ShortURL
	PersistentStorage
}

func NewMemoryStorage(cfg *configs.Config) (*MemoryStorage, error) {

	ps, err := NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}
	ms := &MemoryStorage{
		URLStore:          make(map[string]entity.ShortURL),
		PersistentStorage: ps,
	}

	if err := ps.LoadData(ms.URLStore); err != nil {
		return nil, err
	}

	return ms, err

}
func (ms *MemoryStorage) GetAllURLS(userId string) []entity.ShortURL {
	ms.Lock()
	defer ms.Unlock()
	var urls []entity.ShortURL
	for _, v := range ms.URLStore {
		if v.UserID == userId {
			urls = append(urls, v)
		}
	}
	return urls
}

func (ms *MemoryStorage) GetShortURL(id, uId string) (sURL entity.ShortURL, er error) {
	ms.Lock()
	defer ms.Unlock()

	s, ok := ms.URLStore[id]
	if ok && ms.URLStore[id].UserID == uId {

		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}

func (ms *MemoryStorage) SaveURL(input, userId string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	url := entity.ShortURL{
		ID:     genetareID(),
		Origin: input,
		UserID: userId,
	}

	ms.URLStore[url.ID] = url

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

func (ms *MemoryStorage) Close() error {

	if err := ms.PersistentStorage.SaveData(ms.URLStore); err != nil {
		return err
	}

	if err := ms.PersistentStorage.Close(); err != nil {
		return err
	}
	return nil
}
