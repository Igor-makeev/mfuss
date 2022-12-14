package repositories

import (
	"mfuss/configs"
	"mfuss/internal/entity"
)

type URLStorage interface {
	SaveURL(input string) (string, error)
	GetShortURL(id string) (sURL entity.ShortURL, er error)
}

type PersistentStorage interface {
	SaveData(ms map[string]entity.ShortURL) error
	LoadData(ms map[string]entity.ShortURL) error
}

type Repository struct {
	URLStorage
	PersistentStorage
	Config configs.Config
}

func NewRepository(urls URLStorage, ps PersistentStorage, cfg *configs.Config) *Repository {
	return &Repository{
		URLStorage:        urls,
		PersistentStorage: ps,
		Config:            *cfg,
	}
}
