// Package queries is a struct that contains the queries for the repository.
// It is used to avoid duplicate code in the repository.
package queries

import (
	"context"
	"database/sql"
)

// Queries is a struct that contains the queries for the repository.
type Queries struct {
	DB *sql.DB
}

// NewQueries returns a new instance of the queries.
func NewQueries(db *sql.DB) *Queries {
	return &Queries{
		DB: db,
	}
}

// ExecTx executes a transaction and returns an error if any.
func (r *Queries) ExecTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint:errcheck
	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
