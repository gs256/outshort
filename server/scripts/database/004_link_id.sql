ALTER TABLE links ADD COLUMN uid TEXT;

UPDATE links SET uid = lower(hex(randomblob(16))) WHERE uid IS NULL;

CREATE TABLE links_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    alias TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    owner_id INTEGER,
    name TEXT NOT NULL DEFAULT '',
    lifetime_sec INTEGER NOT NULL DEFAULT 0 CHECK (lifetime_sec >= 0),
    uid TEXT UNIQUE NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO links_new (id, alias, original_url, created_at, owner_id, name, lifetime_sec, uid)
SELECT id, alias, original_url, created_at, owner_id, name, lifetime_sec, uid FROM links;

DROP TABLE links;

ALTER TABLE links_new RENAME TO links;
