package repositories

import (
	"context"
	"mfuss/configs"
	"mfuss/internal/entity"
	errorsEntity "mfuss/internal/entity/errors"
	"mfuss/internal/utilits"
	"mfuss/schema"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// Тип
type PostgresStorage struct {
	DB  *pgx.Conn
	cfg configs.Config
	sync.Mutex
}

// Конструктор
func NewPostgresStorage(cfg *configs.Config, conn *pgx.Conn) *PostgresStorage {

	conn.Exec(context.Background(), schema.Schema)
	conn.Exec(context.Background(), schema.Index)

	ps := &PostgresStorage{
		DB:  conn,
		cfg: *cfg,
	}
	return ps
}

// Клиент постгресса
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

// получить все ссылки
func (ps *PostgresStorage) GetAllURLs(ctx context.Context, userID string) []entity.ShortURL {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	urls := make([]entity.ShortURL, 0)

	rows, err := ps.DB.Query(ctx, `select id,result,origin,user_id from url_store where user_id=$1 ;`, userID)
	if err != nil {
		return nil
	}

	for rows.Next() {
		//Алокация
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

// получить ссылку из хранилища
func (ps *PostgresStorage) GetShortURL(ctx context.Context, id, userID string) (sURL entity.ShortURL, er error) {
	ps.Lock()
	defer ps.Unlock()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	//Алокация
	var url entity.ShortURL
	if err := ps.DB.QueryRow(ctx, `select id,result,origin,user_id, Is_deleted from url_store where id=$1 ;`, id).Scan(&url.ID, &url.ResultURL, &url.Origin, &url.UserID, &url.IsDeleted); err != nil {
		return entity.ShortURL{}, err
	}

	return url, nil
}

// сохранить ссылку
func (ps *PostgresStorage) SaveURL(ctx context.Context, input, userID string) (string, error) {

	ps.Lock()
	defer ps.Unlock()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	//Алокация
	var url entity.ShortURL
	id := utilits.GenetareID()
	res := ps.cfg.BaseURL + "/" + id
	if err := ps.DB.QueryRow(ctx, `insert into url_store(id, result,origin,user_id,Is_deleted) values ($1, $2,$3,$4,$5) on conflict (origin) do update set origin =EXCLUDED.origin, Is_deleted=EXCLUDED.Is_deleted returning *;`, id, res, input, userID, false).Scan(&url.ID, &url.ResultURL, &url.Origin, &url.UserID, &url.IsDeleted); err != nil {

		return "", err
	}

	if url.ID != id {
		return url.ResultURL, errorsEntity.URLConflict{Str: url.Origin}
	}

	return url.ResultURL, nil

}

// закрыть хранилище
func (ps *PostgresStorage) Close(ctx context.Context) error {

	if err := ps.DB.Close(ctx); err != nil {
		return err
	}

	return nil
}

// сохранение батчем
func (ps *PostgresStorage) MultipleShort(ctx context.Context, input []entity.URLBatchInput, userID string) ([]entity.URLBatchResponse, error) {
	//Аллокация
	var resOutput entity.URLBatchResponse
	//Аллокация
	var responseBatch []entity.URLBatchResponse

	for _, v := range input {

		res, err := ps.SaveURL(ctx, v.URL, userID)
		if err != nil {
			return nil, err
		}
		resOutput.CorrelID = v.CorrelID
		resOutput.URL = res
		responseBatch = append(responseBatch, resOutput)

	}

	return responseBatch, nil

}

// проверка соединения с бд
func (ps *PostgresStorage) Ping(ctx context.Context) error {
	err := ps.DB.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil

}

// пометить ссылки как удаленные
func (ps *PostgresStorage) MarkAsDeleted(ctx context.Context, arr []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	_, err := ps.DB.Exec(ctx, "UPDATE url_store SET Is_deleted = true WHERE ID = ANY ($1) and is_deleted <> true", arr)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) GetStats(ctx context.Context) (entity.Stats, error) {
	ps.Lock()
	defer ps.Unlock()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	//Алокация
	var stats entity.Stats
	if err := ps.DB.QueryRow(ctx, `SELECT COUNT(id) AS links_count, COUNT(DISTINCT User_ID) AS users_count FROM links`).Scan(&stats.URLs, &stats.Users); err != nil {
		return stats, err
	}

	return stats, nil
}
