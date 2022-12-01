package repositories

import "mfuss/internal/entity"

type URLStorage interface {
	SaveURL(input string) string
	GetShortURL(id int) (sURL entity.ShortURL, er error)
}

type Repositories struct {
	URLStorage
}

func NewRepository(store URLStorage) *Repositories {
	return &Repositories{
		URLStorage: store,
	}
}
