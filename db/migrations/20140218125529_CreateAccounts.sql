-- +goose Up
CREATE TABLE accounts (
  id          SERIAL    PRIMARY KEY,
  created_at  timestamp NOT NULL,
  email       text      NOT NULL UNIQUE,
  public_key  text      NOT NULL UNIQUE,
  private_key text      NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE accounts;
