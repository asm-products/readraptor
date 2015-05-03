-- +goose Up
ALTER TABLE read_receipts
  ADD COLUMN first_read_at   timestamp;
ALTER TABLE read_receipts
  ADD COLUMN read_count   integer default(0);

-- +goose Down
ALTER TABLE read_receipts
  DROP COLUMN first_read_at;
ALTER TABLE read_receipts
  DROP COLUMN read_count;
