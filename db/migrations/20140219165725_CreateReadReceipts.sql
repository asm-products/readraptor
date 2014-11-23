-- +goose Up
CREATE TABLE read_receipts (
  id              SERIAL    PRIMARY KEY,
  article_id int       NOT NULL,
  reader_id       int       NOT NULL,
  created_at      timestamp NOT NULL,

  CONSTRAINT fk_read_receipts_articles FOREIGN KEY (article_id) REFERENCES articles (id),
  CONSTRAINT fk_read_receipts_readers FOREIGN KEY (reader_id) REFERENCES readers (id),
  CONSTRAINT uq_article_reader UNIQUE(article_id, reader_id)
);

-- +goose Down
DROP TABLE read_receipts;
