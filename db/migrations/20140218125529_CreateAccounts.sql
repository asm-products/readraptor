-- +goose Up
CREATE TABLE accounts (
  id        SERIAL PRIMARY KEY,
  username  text NOT NULL UNIQUE,
  api_key   text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE accounts;

