package repositories

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"
	"strings"

	"github.com/jackc/pgx/v5"

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
	DB     *pgx.Conn
}

func NewRepository(cfg *configs.Config) (*Repository, error) {

	ms, err := NewMemoryStorage(cfg)
	if err != nil {
		return nil, err
	}

	addrCut := strings.TrimPrefix(cfg.DBDSN, "***")
	logrus.Println(addrCut)

	conn, err := pgx.Connect(context.Background(), addrCut)
	if err != nil {
		logrus.Printf("Unable to connect to database: %v\n", err)

	}

	return &Repository{
		URLStorage: ms,
		Config:     *cfg,
		DB:         conn,
	}, nil
}

func (rep *Repository) Close() error {
	if err := rep.URLStorage.Close(); err != nil {
		return err
	}
	if rep.DB != nil {
		err := rep.DB.Close(context.Background())
		if err != nil {
			return err
		}

	}
	return nil
}
