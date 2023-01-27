package repositories

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/utilits"
	"mfuss/schema"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	DB  *pgx.Conn
	cfg configs.Config
	sync.Mutex
}

func NewPostgresStorage(cfg *configs.Config, conn *pgx.Conn) *PostgresStorage {

	conn.Exec(context.Background(), schema.Schema)
	conn.Exec(context.Background(), schema.Index)

	ps := &PostgresStorage{
		DB:  conn,
		cfg: *cfg,
	}
	return ps
}

func NewPostgresClient(cfg *configs.Config) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := pgx.Connect(ctx, cfg.DBDSN)
	if err != nil {
		logrus.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, err
}

func (ps *PostgresStorage) GetAllURLs(userID string, ctx context.Context) []entity.ShortURL {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	urls := make([]entity.ShortURL, 0)

	rows, err := ps.DB.Query(ctx, `select id,result,origin,user_id from url_store where user_id=$1 ;`, userID)
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

func (ps *PostgresStorage) GetShortURL(id, userID string, ctx context.Context) (sURL entity.ShortURL, er error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	var url entity.ShortURL
	if err := ps.DB.QueryRow(ctx, `select id,result,origin,user_id, Is_deleted from url_store where id=$1 ;`, id).Scan(&url.ID, &url.ResultURL, &url.Origin, &url.UserID, &url.IsDeleted); err != nil {
		return entity.ShortURL{}, err
	}

	return url, nil
}

func (ps *PostgresStorage) SaveURL(input, userID string, ctx context.Context) (string, error) {
	ps.Lock()
	defer ps.Unlock()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	var url entity.ShortURL
	id := utilits.GenetareID()
	res := ps.cfg.BaseURL + "/" + id
	if err := ps.DB.QueryRow(ctx, `insert into url_store(id, result,origin,user_id,Is_deleted) values ($1, $2,$3,$4,$5) on conflict (origin) do update set origin =EXCLUDED.origin, Is_deleted=EXCLUDED.Is_deleted returning *;`, id, res, input, userID, false).Scan(&url.ID, &url.ResultURL, &url.Origin, &url.UserID, &url.IsDeleted); err != nil {

		return "", err
	}

	if url.ID != id {
		return url.ResultURL, utilits.URLConflict{Str: url.Origin}
	}

	return url.ResultURL, nil

}

func (ps *PostgresStorage) Close(ctx context.Context) error {
	if _, err := ps.DB.Exec(ctx, "Drop table url_store;"); err != nil {
		return err
	}

	if err := ps.DB.Close(ctx); err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStorage) MultipleShort(input []entity.URLBatchInput, userID string, ctx context.Context) ([]entity.URLBatchResponse, error) {
	var resOutput entity.URLBatchResponse
	var responseBatch []entity.URLBatchResponse

	for _, v := range input {
		res, err := ps.SaveURL(v.URL, userID, ctx)
		if err != nil {
			return nil, err
		}
		resOutput.CorrelID = v.CorrelID
		resOutput.URL = res
		responseBatch = append(responseBatch, resOutput)

	}

	return responseBatch, nil

}

func (ps *PostgresStorage) Ping(ctx context.Context) error {
	err := ps.DB.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil

}

func (ps *PostgresStorage) MarkAsDeleted(arr []string) error {

	_, err := ps.DB.Exec(context.Background(), "UPDATE url_store SET Is_deleted = true WHERE ID = ANY ($1) and is_deleted <> true", arr)
	if err != nil {
		return err
	}
	return nil
}
