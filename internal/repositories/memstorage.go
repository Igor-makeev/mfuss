package repositories

import (
	"context"
	"errors"
	"fmt"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/utilits"
	"sync"

	"github.com/sirupsen/logrus"
)

type Dumper interface {
	SaveData(ms map[string]*entity.ShortURL) error
	LoadData(ms map[string]*entity.ShortURL) error
	Close() error
}

type MemoryStorage struct {
	sync.Mutex
	URLStore map[string]*entity.ShortURL
	cfg      configs.Config
	Dumper
}

func NewMemoryStorage(cfg *configs.Config, dumper Dumper) *MemoryStorage {

	return &MemoryStorage{
		URLStore: make(map[string]*entity.ShortURL),
		Dumper:   dumper,
		cfg:      *cfg,
	}

}

func (ms *MemoryStorage) LoadFromDump() error {
	return ms.Dumper.LoadData(ms.URLStore)
}

func (ms *MemoryStorage) GetAllURLs(userID string, ctx context.Context) []entity.ShortURL {

	var urls []entity.ShortURL
	for _, v := range ms.URLStore {
		if v.UserID == userID {
			urls = append(urls, *v)
		}
	}
	return urls
}

func (ms *MemoryStorage) GetShortURL(id, userID string, ctx context.Context) (sURL entity.ShortURL, er error) {

	s, ok := ms.URLStore[id]
	if ok {

		return *s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}

func (ms *MemoryStorage) SaveURL(input, userID string, ctx context.Context) (string, error) {
	ms.Lock()
	defer ms.Unlock()
	for _, value := range ms.URLStore {
		if value.Origin == input {
			return value.ResultURL, utilits.URLConflict{Str: value.Origin}
		}
	}

	url := &entity.ShortURL{
		ID:     utilits.GenetareID(),
		Origin: input,
		UserID: userID,
	}

	url.ResultURL = ms.cfg.BaseURL + "/" + url.ID
	ms.URLStore[url.ID] = url

	return url.ResultURL, nil
}

func (ms *MemoryStorage) Close(ctx context.Context) error {

	if err := ms.Dumper.SaveData(ms.URLStore); err != nil {
		return err
	}

	return ms.Dumper.Close()
}

func (ms *MemoryStorage) MultipleShort(input []entity.URLBatchInput, userID string, ctx context.Context) ([]entity.URLBatchResponse, error) {
	var resOutput entity.URLBatchResponse
	var responseBatch []entity.URLBatchResponse

	for _, v := range input {
		res, err := ms.SaveURL(v.URL, userID, ctx)
		if err != nil {
			return nil, err
		}
		resOutput.CorrelID = v.CorrelID
		resOutput.URL = res
		responseBatch = append(responseBatch, resOutput)

	}

	return responseBatch, nil

}

func (ms *MemoryStorage) Ping(ctx context.Context) error {

	return errors.New("no db connection")

}

func (ms *MemoryStorage) MarkAsDeleted(arr []string, ctx context.Context) error {

	for _, val := range arr {
		ms.setDeletFlag(val)
		logrus.Printf("in markasdeleted %v", ms.URLStore[val])
	}

	return nil
}

func (ms *MemoryStorage) setDeletFlag(ID string) {

	for i, v := range ms.URLStore {
		if i == ID {
			v.SetDeleteFlag()

		}
		logrus.Printf("in setdeleteflag %v", v)
	}

}
