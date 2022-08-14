package repository

import (
	"database/sql"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/auth"
)

type Auth interface {
	CreateUser(userInput domain.SignUpInput) (int, error)
	GetUser(email, password string) (domain.User, error)
}

type Repository struct {
	Auth
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		return nil
	}

	return &Repository{
		Auth: auth.NewRepoAuth(db),
	}
}
