package repositories

import (
	"context"
	"mfuss/configs"

	"mfuss/internal/entity"
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

func NewRepository(cfg *configs.Config, urlstorager URLStorager) (*Repository, error) {

	return &Repository{
		URLStorager: urlstorager,
		Config:      cfg,
	}, nil

}

func (rep *Repository) Close(ctx context.Context) error {

	if err := rep.URLStorager.Close(ctx); err != nil {
		return err
	}

	return nil
}
