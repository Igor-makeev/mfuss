package schema

var Schema = `
CREATE TABLE url_store (
    ID text,
    Result text,
    Origin text,
	User_ID text
);`

var Index = `
CREATE UNIQUE INDEX url_store_index_unique
  ON url_store
  USING btree(Origin);
`
