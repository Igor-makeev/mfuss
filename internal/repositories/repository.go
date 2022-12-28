package repositories

import (
	"mfuss/configs"
	"mfuss/internal/entity"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
	DB     *sqlx.DB
}

func NewRepository(cfg *configs.Config) (*Repository, error) {

	ms, err := NewMemoryStorage(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("postgres", cfg.DBDSN)
	if err != nil {
		logrus.Println("Error connecting to database")
	}

	return &Repository{
		URLStorage: ms,
		Config:     *cfg,
		DB:         db,
	}, nil
}
