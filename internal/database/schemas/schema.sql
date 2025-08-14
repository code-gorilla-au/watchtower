CREATE TABLE IF NOT EXISTS organisations
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    friendly_name TEXT    NOT NULL,
    namespace     TEXT    NOT NULL UNIQUE,
    default_org   BOOLEAN NOT NULL DEFAULT 0,
    token         TEXT    NOT NULL,
    created_at    INTEGER NOT NULL,
    updated_at    INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS products
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT UNIQUE NOT NULL,
    tags       TEXT, -- JSON array of strings
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS product_organisations
(
    product_id      INTEGER REFERENCES products (id),
    organisation_id INTEGER REFERENCES organisations (id),
    PRIMARY KEY (product_id, organisation_id)
);

CREATE TABLE IF NOT EXISTS repositories
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT UNIQUE NOT NULL,
    url        TEXT NOT NULL,
    topic      TEXT NOT NULL,
    owner      TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS securities
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id     TEXT UNIQUE NOT NULL,
    repository_name TEXT NOT NULL,
    package_name    TEXT NOT NULL,
    state           TEXT NOT NULL,
    severity        TEXT NOT NULL,
    patched_version TEXT NOT NULL,
    created_at      INTEGER NOT NULL,
    updated_at      INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id     TEXT UNIQUE NOT NULL,
    title           TEXT NOT NULL,
    repository_name TEXT NOT NULL,
    url             TEXT NOT NULL,
    state           TEXT NOT NULL,
    author          TEXT NOT NULL,
    merged_at       INTEGER NOT NULL,
    created_at      INTEGER NOT NULL,
    updated_at      INTEGER NOT NULL
);