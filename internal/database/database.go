package database

import (
	"database/sql"
	_ "embed"
	"path"

	_ "modernc.org/sqlite"
)

//go:embed schemas/schema.sql
var ddl string

const dbName = "watchtower.db"

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Init() error {
	_, err := m.db.Exec(ddl)

	return err
}

// NewDBFromProvider initializes a new Queries instance using a SQLite database at the given file path.
// Returns a pointer to Queries and an error if the database connection cannot be established.
func NewDBFromProvider(filePath string) (*Queries, *sql.DB, error) {
	dbFile := resolveDBPath(filePath)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, nil, err
	}

	return New(db), db, nil
}

func resolveDBPath(filePath string) string {
	if filePath == ":memory:" {
		return filePath
	}

	return path.Join(filePath, dbName, "?_busy_timeout=5000")
}
