package repositories

import (
	"mfuss/configs"
	"mfuss/internal/entity"
	"strings"

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

	addrCut := strings.TrimPrefix(cfg.DBDSN, "***")
	logrus.Println(addrCut)
	db, err := sqlx.Open("postgres", addrCut)
	if err != nil {
		logrus.Println("Error connecting to database")
	}

	return &Repository{
		URLStorage: ms,
		Config:     *cfg,
		DB:         db,
	}, nil
}

func (rep *Repository) Close() error {
	if err := rep.URLStorage.Close(); err != nil {
		return err
	}
	if rep.DB != nil {
		err := rep.DB.Close()
		if err != nil {
			return err
		}

	}
	return nil
}
