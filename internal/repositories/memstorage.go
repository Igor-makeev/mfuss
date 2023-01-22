package repositories

import (
	"errors"
	"fmt"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/utilits"
	"sync"
)

type Dumper interface {
	SaveData(ms map[string]entity.ShortURL) error
	LoadData(ms map[string]entity.ShortURL) error
	Close() error
}

type MemoryStorage struct {
	sync.Mutex
	URLStore map[string]entity.ShortURL
	cfg      configs.Config
	Dumper
}

func NewMemoryStorage(cfg *configs.Config, d Dumper) (*MemoryStorage, error) {

	ms := &MemoryStorage{
		URLStore: make(map[string]entity.ShortURL),
		Dumper:   d,
		cfg:      *cfg,
	}

	err := d.LoadData(ms.URLStore)
	if err != nil {
		return nil, err
	}

	return ms, err

}
func (ms *MemoryStorage) GetAllURLS(userID string) []entity.ShortURL {
	ms.Lock()
	defer ms.Unlock()
	var urls []entity.ShortURL
	for _, v := range ms.URLStore {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}
	return urls
}

func (ms *MemoryStorage) GetShortURL(id, userID string) (sURL entity.ShortURL, er error) {
	ms.Lock()
	defer ms.Unlock()

	s, ok := ms.URLStore[id]
	if ok {

		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}

func (ms *MemoryStorage) SaveURL(input, userID string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	url := entity.ShortURL{
		ID:     utilits.GenetareID(),
		Origin: input,
		UserID: userID,
	}
	url.ResultURL = ms.cfg.BaseURL + "/" + url.ID
	ms.URLStore[url.ID] = url

	return url.ResultURL, nil
}

func (ms *MemoryStorage) Close() error {

	if err := ms.Dumper.SaveData(ms.URLStore); err != nil {
		return err
	}

	if err := ms.Dumper.Close(); err != nil {
		return err
	}
	return nil
}

func (ms *MemoryStorage) MultipleShort(input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	var resOutput entity.URLBatchResponse
	var responseBatch []entity.URLBatchResponse

	for _, v := range input {
		res, err := ms.SaveURL(v.URL, userID)
		if err != nil {
			return nil, err
		}
		resOutput.CorrelID = v.CorrelID
		resOutput.URL = res
		responseBatch = append(responseBatch, resOutput)

	}

	return responseBatch, nil

}

func (ms *MemoryStorage) Ping() error {

	return errors.New("no db connection")

}

func (ms *MemoryStorage) MarkAsDeleted(arr []string, id string) error {

	return nil

}
