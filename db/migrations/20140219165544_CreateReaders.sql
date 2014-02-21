-- +goose Up
CREATE TABLE readers (
  id          SERIAL    PRIMARY KEY,
  account_id  int       NOT NULL,
  created_at  timestamp NOT NULL,
  distinct_id text      NOT NULL UNIQUE,

  CONSTRAINT fk_readers_account FOREIGN KEY (account_id) REFERENCES accounts (id)
);

-- +goose Down
DROP TABLE readers;

