// Package postgres implements database driver for postgres
package postgres

import (
	"database/sql"
	"github.com/popeskul/qna-go/internal/db"
)

// NewPostgresConnection creates new postgres connection and returns *sql.DB and error.
func NewPostgresConnection(cfg db.ConfigDB) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.String())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(300)
	db.SetMaxOpenConns(60)

	return db, nil
}
