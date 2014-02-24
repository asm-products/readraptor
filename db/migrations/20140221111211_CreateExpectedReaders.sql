-- +goose Up
CREATE TABLE expected_readers (
  id              SERIAL    PRIMARY KEY,
  article_id int       NOT NULL,
  reader_id       int       NOT NULL,
  created_at      timestamp NOT NULL,

  CONSTRAINT fk_expected_readers_articles FOREIGN KEY (article_id) REFERENCES articles (id),
  CONSTRAINT fk_expected_readers_readers FOREIGN KEY (reader_id) REFERENCES readers (id),
  CONSTRAINT uq_expected_readers_article_reader UNIQUE(article_id, reader_id)
);

-- +goose Down
DROP TABLE expected_readers;
