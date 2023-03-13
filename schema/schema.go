package schema

// схема таблицы хранилища ссылок
var Schema = `
create table if not exists url_store (
    ID text,
    Result text,
    Origin text,
	User_ID text,
  Is_deleted bool
);`

// индекс уникальных ссылок
var Index = `
CREATE UNIQUE INDEX if not exists url_store_index_unique
  ON url_store
  USING btree(Origin);
`
