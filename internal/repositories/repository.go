package repositories

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"

	"github.com/jackc/pgx/v5"
)

type URLStorager interface {
	SaveURL(input, userID string) (string, error)
	GetAllURLS(userID string) []entity.ShortURL
	GetShortURL(id, userID string) (sURL entity.ShortURL, er error)
	MultipleShort(input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error)
	Close() error
}

type Repository struct {
	URLStorager
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
			URLStorager: ms,
			Config:      *cfg,
			DB:          nil,
		}, nil
	}

	ps, err := NewPostgresStorage(cfg)
	if err != nil {
		return nil, err
	}
	return &Repository{
		URLStorager: ps,
		Config:      *cfg,
		DB:          ps.DB,
	}, nil

}

func (rep *Repository) Close() error {

	if err := rep.URLStorager.Close(); err != nil {
		return err
	}

	if rep.DB != nil {
		if err := rep.DB.Close(context.Background()); err != nil {
			return err
		}
	}

	return nil
}
