// Package database provides functionality for interacting with the SQLite database, including migrations and query initialization.
package database

import (
	"database/sql"
	_ "embed"
	"path"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed schemas/schema.sql
var ddl string

const dbName = "watchtower.db"

// Migrator handles database schema initialization and migrations.
type Migrator struct {
	db *sql.DB
}

// NewMigrator creates a new Migrator instance with the provided database connection.
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Init initializes the database schema by executing the DDL.
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

// resolveDBPath resolves the database file path, handling the special ":memory:" case.
func resolveDBPath(filePath string) string {
	if filePath == ":memory:" {
		return filePath
	}

	return path.Join(filePath, dbName, "?_busy_timeout=5000")
}

// IsErrUniqueConstraint checks if the error is a SQLite unique constraint violation.
func IsErrUniqueConstraint(err error) bool {
	return err != nil && strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed")
}
