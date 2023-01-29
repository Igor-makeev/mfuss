package repositories

import (
	"context"
	"mfuss/configs"

	"mfuss/internal/entity"

	"github.com/sirupsen/logrus"
)

type URLStorager interface {
	SaveURL(input, userID string, ctx context.Context) (string, error)
	GetAllURLs(userID string, ctx context.Context) []entity.ShortURL
	GetShortURL(id, userID string, ctx context.Context) (sURL entity.ShortURL, er error)
	MultipleShort(input []entity.URLBatchInput, userID string, ctx context.Context) ([]entity.URLBatchResponse, error)
	MarkAsDeleted(arr []string) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

type Repository struct {
	URLStorager
	Config *configs.Config
}

func NewRepository(cfg *configs.Config) (*Repository, error) {

	var urlstorage URLStorager
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

func (rep *Repository) Close(ctx context.Context) error {

	if err := rep.URLStorager.Close(ctx); err != nil {
		return err
	}

	return nil
}
