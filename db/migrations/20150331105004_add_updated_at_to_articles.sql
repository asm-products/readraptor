-- +goose Up
ALTER TABLE articles
  ADD COLUMN updated_at   timestamp;

update articles set updated_at=created_at;

-- +goose Down
ALTER TABLE articles
  DROP COLUMN updated_at;
