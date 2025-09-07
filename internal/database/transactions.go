package database

import (
	"context"
	"database/sql"

	"watchtower/internal/logging"
)

type TxnFn func(tx *sql.Tx) error

//go:generate moq -rm -stub -out mocks.gen.go . DBTxner DBBeginner

type DBBeginner interface {
	Begin() (*sql.Tx, error)
}

type DBTxner interface {
	DBBeginner
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	DBTX
}

// WithTxnContext executes a function within a transaction context.
// If the function returns an error, the transaction is rolled back.
func WithTxnContext(ctx context.Context, db DBBeginner, txnFn TxnFn) error {
	log := logging.FromContext(ctx)

	tx, err := db.Begin()
	if err != nil {
		log.Error("database transaction failed", "error", err)
		return err
	}

	if err = txnFn(tx); err != nil {
		log.Error("transaction func failed", "error", err)

		if rbErr := tx.Rollback(); rbErr != nil {
			log.Error("transaction rollback failed", "error", err)
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
