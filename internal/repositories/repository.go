package repositories

import "mfuss/internal/entity"

type URLStorage interface {
	SaveURL(input string) (string, error)
	GetShortURL(id string) (sURL entity.ShortURL, er error)
}
