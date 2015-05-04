-- +goose Up
ALTER TABLE articles
  ADD COLUMN total_read_count integer default(0);
ALTER TABLE articles
  ADD COLUMN unique_read_count integer default(0);

-- +goose Down
ALTER TABLE read_receipts
  DROP COLUMN total_read_count;
ALTER TABLE read_receipts
  DROP COLUMN unique_read_count;
