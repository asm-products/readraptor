-- +goose Up
ALTER TABLE accounts
  ADD COLUMN customer_id   text;

-- +goose Down
ALTER TABLE accounts
  DROP COLUMN customer_id;

