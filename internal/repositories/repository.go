package repositories

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"

	"github.com/jackc/pgx/v5"
)

type URLStorage interface {
	SaveURL(input, userID string) (string, error)
	GetAllURLS(userID string) []entity.ShortURL
	GetShortURL(id, userID string) (sURL entity.ShortURL, er error)
	Close() error
}

type Repository struct {
	URLStorage
	Config configs.Config
	DB     *pgx.Conn
}

func NewRepository(cfg *configs.Config) (*Repository, error) {

	if cfg.DBDSN == "" {
		ms, err := NewMemoryStorage(cfg)
		if err != nil {
			return nil, err
		}
		return &Repository{
			URLStorage: ms,
			Config:     *cfg,
			DB:         nil,
		}, nil
	}

	ps, err := NewPostgresStorage(cfg)
	if err != nil {
		return nil, err
	}
	return &Repository{
		URLStorage: ps,
		Config:     *cfg,
		DB:         ps.DB,
	}, nil

}

func (rep *Repository) Close() error {
	if rep.URLStorage != nil {
		if err := rep.URLStorage.Close(); err != nil {
			return err
		}
	}
	if rep.DB != nil {
		if err := rep.DB.Close(context.Background()); err != nil {
			return err
		}
	}

	return nil
}
