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

func NewMigrator(filePath string) (*Migrator, error) {
	dbFile := path.Join(filePath, dbName)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	return &Migrator{db: db}, nil
}

func (m *Migrator) Init() error {
	_, err := m.db.Exec(ddl)
	return err
}

// NewDBFromProvider initializes a new Queries instance using a SQLite database at the given file path.
// Returns a pointer to Queries and an error if the database connection cannot be established.
func NewDBFromProvider(filePath string) (*Queries, error) {
	dbFile := path.Join(filePath, dbName)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	return New(db), nil
}
