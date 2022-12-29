package repositories

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/utilits"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	DB  *pgx.Conn
	cfg configs.Config
	sync.Mutex
}

var schema = `
CREATE TABLE url_store (
    ID text,
    Result text,
    Origin text,
	User_ID text
);`

func NewPostgresStorage(cfg *configs.Config) (*PostgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DBDSN)
	if err != nil {
		logrus.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	conn.Exec(context.Background(), schema)

	ps := &PostgresStorage{
		DB:  conn,
		cfg: *cfg,
	}
	return ps, err
}

func (ps *PostgresStorage) GetAllURLS(userID string) []entity.ShortURL {
	ps.Lock()
	defer ps.Unlock()
	urls := make([]entity.ShortURL, 0)

	rows, err := ps.DB.Query(context.Background(), `select id,result,origin,user_id from url_store where user_id=$1;`, userID)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var url entity.ShortURL
		err = rows.Scan(&url.ID, &url.ResultURL, &url.Origin, &url.UserID)
		if err != nil {
			return nil
		}
		urls = append(urls, url)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}

	return urls
}

func (ps *PostgresStorage) GetShortURL(id, userID string) (sURL entity.ShortURL, er error) {
	ps.Lock()
	defer ps.Unlock()
	var url entity.ShortURL
	if err := ps.DB.QueryRow(context.Background(), `select id,result,origin,user_id from url_store where id=$1 and user_id=$2;`, id, userID).Scan(&url.ID, &url.ResultURL, &url.Origin, &url.UserID); err != nil {
		return entity.ShortURL{}, err
	}

	return url, nil
}

func (ps *PostgresStorage) SaveURL(input, userID string) (string, error) {
	ps.Lock()
	defer ps.Unlock()

	id := utilits.GenetareID()
	res := ps.cfg.BaseURL + "/" + id
	if _, err := ps.DB.Exec(context.Background(), `insert into url_store(id, result,origin,user_id) values ($1, $2,$3,$4);`, id, res, input, userID); err != nil {
		return "", err
	}
	return res, nil

}

func (ps *PostgresStorage) Close() error {
	// if _, err := ps.DB.Exec(context.Background(), "Drop table url_store;"); err != nil {
	// 	return err
	// }
	// TODO fix this moment with drop
	if err := ps.DB.Close(context.Background()); err != nil {
		return err
	}

	return nil
}
