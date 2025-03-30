-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    alias TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
