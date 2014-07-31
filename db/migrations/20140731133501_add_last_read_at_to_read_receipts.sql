-- +goose Up
ALTER TABLE read_receipts
  ADD COLUMN last_read_at   timestamp;

update read_receipts set last_read_at=created_at;

-- +goose Down
ALTER TABLE read_receipts
  DROP COLUMN last_read_at;
