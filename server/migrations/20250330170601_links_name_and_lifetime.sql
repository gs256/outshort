-- +goose Up
-- +goose StatementBegin
ALTER TABLE links
ADD COLUMN name TEXT NOT NULL DEFAULT '';

ALTER TABLE links
ADD COLUMN lifetime_sec INTEGER NOT NULL DEFAULT 0 CHECK (lifetime_sec >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links DROP COLUMN name;
ALTER TABLE links DROP COLUMN lifetime_sec;
-- +goose StatementEnd
