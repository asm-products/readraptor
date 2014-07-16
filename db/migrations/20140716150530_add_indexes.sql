-- +goose Up
CREATE INDEX expected_readers_reader_id ON expected_readers (reader_id)
CREATE INDEX read_receipts_reader_id ON read_receipts (reader_id)

-- +goose Down
DROP INDEX expected_readers_reader_id
DROP INDEX read_receipts_reader_id
