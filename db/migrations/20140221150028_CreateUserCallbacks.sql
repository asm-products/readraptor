-- +goose Up
CREATE TABLE user_callbacks (
  id              SERIAL    PRIMARY KEY,
  content_item_id int       NOT NULL,
  reader_id       int       NOT NULL,
  created_at      timestamp NOT NULL,

  CONSTRAINT fk_user_callbacks_content_items FOREIGN KEY (content_item_id) REFERENCES content_items (id),
  CONSTRAINT fk_user_callbacks_readers FOREIGN KEY (reader_id) REFERENCES readers (id),
  CONSTRAINT uq_user_callbacks_content_item_reader UNIQUE(content_item_id, reader_id)
);

-- +goose Down
DROP TABLE expected_readers;
