CREATE TABLE IF NOT EXISTS organisations (
                                             id INTEGER PRIMARY KEY AUTOINCREMENT,
                                             friendly_name TEXT NOT NULL,
                                             namespace TEXT NOT NULL UNIQUE,
                                             created_at INTEGER NOT NULL,
                                             updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
                                        id TEXT PRIMARY KEY,
                                        name TEXT UNIQUE,
                                        tags TEXT, -- JSON array of strings
                                        created_at INTEGER,
                                        updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS product_organisations (
                                                     product_id TEXT REFERENCES products(id),
    organisation_id INTEGER REFERENCES organisations(id),
    PRIMARY KEY (product_id, organisation_id)
    );

CREATE TABLE IF NOT EXISTS repositories (
                                            id TEXT PRIMARY KEY,
                                            name TEXT UNIQUE,
                                            url TEXT,
                                            topic TEXT,
                                            owner TEXT,
                                            created_at INTEGER,
                                            updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS securities (
                                          id TEXT PRIMARY KEY,
                                          external_id TEXT UNIQUE,
                                          repository_name TEXT,
                                          package_name TEXT,
                                          state TEXT,
                                          severity TEXT,
                                          patched_version TEXT,
                                          created_at INTEGER,
                                          updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS pull_requests (
                                             id TEXT PRIMARY KEY,
                                             external_id TEXT UNIQUE,
                                             title TEXT,
                                             repository_name TEXT,
                                             url TEXT,
                                             state TEXT,
                                             author TEXT,
                                             merged_at INTEGER,
                                             created_at INTEGER,
                                             updated_at INTEGER
);