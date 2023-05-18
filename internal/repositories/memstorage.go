package repositories

import (
	"context"
	"errors"
	"fmt"
	"mfuss/configs"
	"mfuss/internal/entity"
	errorsEntity "mfuss/internal/entity/errors"
	"mfuss/internal/utilits"
	"sync"
)

// Интерфейс дампер
type Dumper interface {
	SaveData(ms map[string]*entity.ShortURL) error
	LoadData(ms map[string]*entity.ShortURL) error
	Close() error
}

// Тип мемори сторэдж
type MemoryStorage struct {
	sync.Mutex
	URLStore map[string]*entity.ShortURL
	cfg      configs.Config
	Dumper
}

// Конструктор
func NewMemoryStorage(cfg *configs.Config, dumper Dumper) *MemoryStorage {

	return &MemoryStorage{
		URLStore: make(map[string]*entity.ShortURL, 5),
		Dumper:   dumper,
		cfg:      *cfg,
	}

}

// Загружает данные из дампа
func (ms *MemoryStorage) LoadFromDump() error {
	return ms.Dumper.LoadData(ms.URLStore)
}

// Получить все сокращенные ссылки пользователя
func (ms *MemoryStorage) GetAllURLs(ctx context.Context, userID string) []entity.ShortURL {
	//Алокация
	var urls []entity.ShortURL
	for _, v := range ms.URLStore {
		if v.UserID == userID {
			urls = append(urls, *v)
		}
	}
	return urls
}

// Получить сокращенную ссылку
func (ms *MemoryStorage) GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error) {

	s, ok := ms.URLStore[id]
	if ok {

		return *s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%v not found", id)

}

// Сохранить сокращенную ссылку
func (ms *MemoryStorage) SaveURL(ctx context.Context, input, userID string) (string, error) {
	ms.Lock()
	defer ms.Unlock()
	for _, value := range ms.URLStore {
		if value.Origin == input {
			return value.ResultURL, errorsEntity.URLConflict{Str: value.Origin}
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

// закрыть хранилище
func (ms *MemoryStorage) Close(ctx context.Context) error {

	if err := ms.Dumper.SaveData(ms.URLStore); err != nil {
		return err
	}

	return ms.Dumper.Close()
}

// Сохранение батчами
func (ms *MemoryStorage) MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	//Аллокация
	var resOutput entity.URLBatchResponse
	//Аллокация
	var responseBatch []entity.URLBatchResponse

	for _, v := range input {
		res, err := ms.SaveURL(ctx, v.URL, userID)
		if err != nil {
			return nil, err
		}
		resOutput.CorrelID = v.CorrelID
		resOutput.URL = res
		responseBatch = append(responseBatch, resOutput)

	}

	return responseBatch, nil

}

// Заглушка метода проверки связи
func (ms *MemoryStorage) Ping(ctx context.Context) error {

	return errors.New("no db connection")

}

// пометить ссылку как удаленную
func (ms *MemoryStorage) MarkAsDeleted(ctx context.Context, arr []string) error {

	for _, val := range arr {
		ms.setDeletFlag(val)

	}

	return nil
}

// установить флаг удаления
func (ms *MemoryStorage) setDeletFlag(ID string) {

	for i, v := range ms.URLStore {
		if i == ID {
			v.SetDeleteFlag()

		}

	}

}

// получить статы
func (ms *MemoryStorage) GetStats(ctx context.Context) (entity.Stats, error) {
	ms.Lock()
	defer ms.Unlock()
	users := make(map[string]bool)

	for _, v := range ms.URLStore {
		users[v.UserID] = true
	}

	return entity.Stats{
		URLs:  int(len(ms.URLStore)),
		Users: int(len(users)),
	}, nil

}
