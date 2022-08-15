package repository

import (
	"database/sql"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/auth"
	"github.com/popeskul/qna-go/internal/repository/tests"
)

type Auth interface {
	CreateUser(userInput domain.SignUpInput) (int, error)
	GetUser(email, password string) (domain.User, error)
}

type Tests interface {
	CreateTest(userID int, testInput domain.TestInput) (int, error)
	GetTest(testID int) (domain.Test, error)
}

type Repository struct {
	Auth
	Tests
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		return nil
	}

	return &Repository{
		Auth:  auth.NewRepoAuth(db),
		Tests: tests.NewRepoTests(db),
	}
}
