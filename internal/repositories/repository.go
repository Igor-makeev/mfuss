package repositories

import (
	"context"
	"mfuss/configs"

	"mfuss/internal/entity"

	"github.com/sirupsen/logrus"
)

// Интерфейс хранилища ссылок
type URLStorager interface {
	SaveURL(ctx context.Context, input, userID string) (string, error)
	GetAllURLs(ctx context.Context, userID string) []entity.ShortURL
	GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error)
	MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error)
	MarkAsDeleted(ctx context.Context, arr []string) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

// Структура репозитория
type Repository struct {
	URLStorager
	Config *configs.Config
}

// КОнструктор
func NewRepository(cfg *configs.Config) (*Repository, error) {
	//Алокация
	var urlstorage URLStorager
	//Алокация
	var err error

	if cfg.DBDSN == "" {
		urlstorage, err = PrepareMemoryStorage(cfg)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		conn, err := NewPostgresClient(cfg)
		if err != nil {
			logrus.Fatal(err)
		}
		urlstorage = NewPostgresStorage(cfg, conn)
	}

	return &Repository{
		URLStorager: urlstorage,
		Config:      cfg,
	}, nil

}

// Подготовить хранилище данных
func PrepareMemoryStorage(cfg *configs.Config) (*MemoryStorage, error) {

	dump, err := NewDump(cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}

	ms := NewMemoryStorage(cfg, dump)
	if err := ms.LoadFromDump(); err != nil {
		return nil, err
	}

	return ms, nil

}

// закрыть хранилище данных
func (rep *Repository) Close(ctx context.Context) error {

	if err := rep.URLStorager.Close(ctx); err != nil {
		return err
	}

	return nil
}
