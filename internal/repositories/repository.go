package repositories

import (
	"mfuss/configs"
	"mfuss/internal/entity"
)

type URLStorager interface {
	SaveURL(input, userID string) (string, error)
	GetAllURLS(userID string) []entity.ShortURL
	GetShortURL(id, userID string) (sURL entity.ShortURL, er error)
	MultipleShort(input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error)
	MarkAsDeleted(arr []string, id string)
	Ping() error
	Close() error
}

type Repository struct {
	URLStorager
	Config configs.Config
	Buffer *Buffer
}

func NewRepository(cfg *configs.Config, urlstorager URLStorager) (*Repository, error) {

	return &Repository{
		URLStorager: urlstorager,
		Config:      *cfg,
		Buffer:      NewBuffer(),
	}, nil

}

func (rep *Repository) Close() error {

	if err := rep.URLStorager.Close(); err != nil {
		return err
	}

	return nil
}
