package postgres

import (
	"database/sql"
	"github.com/popeskul/qna-go/internal/db"
)

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
