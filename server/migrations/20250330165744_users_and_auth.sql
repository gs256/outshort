-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    username TEXT UNIQUE NOT NULL,
    -- FIXME!
    password TEXT NOT NULL
);

CREATE TABLE auth_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    token TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    lifetime_sec INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE links_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    alias TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    owner_id INTEGER,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO links_new (id, alias, original_url, created_at)
SELECT id, alias, original_url, created_at FROM links;

DROP TABLE links;

ALTER TABLE links_new RENAME TO links;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth_tokens;
DROP TABLE users;

ALTER TABLE links RENAME TO links_old;

CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    alias TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO links (id, alias, original_url, created_at)
SELECT id, alias, original_url, created_at FROM links_old;

DROP TABLE links_old;
-- +goose StatementEnd
