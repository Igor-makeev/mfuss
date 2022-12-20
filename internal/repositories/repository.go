package repositories

import (
	"mfuss/configs"
	"mfuss/internal/entity"
)

type URLStorage interface {
	SaveURL(input string) (string, error)
	GetShortURL(id string) (sURL entity.ShortURL, er error)
	Close() error
}

type PersistentStorage interface {
	SaveData(ms map[string]entity.ShortURL) error
	LoadData(ms map[string]entity.ShortURL) error
	Close() error
}

type Repository struct {
	URLStorage
	Config configs.Config
}

func NewRepository(cfg *configs.Config) (*Repository, error) {

	ms, err := NewMemoryStorage(cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		URLStorage: ms,
		Config:     *cfg,
	}, nil
}
