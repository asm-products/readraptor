-- +goose Up
CREATE TABLE articles (
  id         SERIAL    PRIMARY KEY,
  account_id int       NOT NULL,
  created_at timestamp NOT NULL,
  key        text      NOT NULL UNIQUE,

  CONSTRAINT fk_articles_account FOREIGN KEY (account_id) REFERENCES accounts (id)
);

-- +goose Down
DROP TABLE articles;

