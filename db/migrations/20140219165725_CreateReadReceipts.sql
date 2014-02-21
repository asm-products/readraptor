-- +goose Up
CREATE TABLE read_receipts (
  id              SERIAL    PRIMARY KEY,
  content_item_id int       NOT NULL,
  reader_id       int       NOT NULL,
  created_at      timestamp NOT NULL,

  CONSTRAINT fk_read_receipts_content_items FOREIGN KEY (content_item_id) REFERENCES content_items (id),
  CONSTRAINT fk_read_receipts_readers FOREIGN KEY (reader_id) REFERENCES readers (id),
  CONSTRAINT uq_content_item_reader UNIQUE(content_item_id, reader_id)
);

-- +goose Down
DROP TABLE read_receipts;
