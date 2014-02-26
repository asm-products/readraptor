-- +goose Up
ALTER TABLE accounts
  ADD COLUMN confirmation_token   text,
  ADD COLUMN confirmation_sent_at timestamp,
  ADD COLUMN confirmed_at         timestamp;

-- +goose Down
ALTER TABLE accounts
  DROP COLUMN confirmation_token,
  DROP COLUMN confirmation_sent_at,
  DROP COLUMN confirmed_at,

