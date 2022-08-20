// Package repository is a struct that contains the repository.
// This place define interface for the repository: Auth, Tests.
package repository

import (
	"database/sql"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/auth"
	"github.com/popeskul/qna-go/internal/repository/tests"
)

// Auth interface is implemented by the auth repository.
type Auth interface {
	CreateUser(userInput domain.SignUpInput) (int, error)
	GetUser(email, password string) (domain.User, error)
	DeleteUserById(userID int) error
}

// Tests interface is implemented by the test repository.
type Tests interface {
	CreateTest(userID int, testInput domain.TestInput) (int, error)
	GetTest(testID int) (domain.Test, error)
	UpdateTestById(testID int, testInput domain.TestInput) error
	DeleteTestById(testID int) error
}

// Repository is the composite of all repositories.
type Repository struct {
	Auth
	Tests
}

// NewRepository returns a new instance of the repository.
func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		return nil
	}

	return &Repository{
		Auth:  auth.NewRepoAuth(db),
		Tests: tests.NewRepoTests(db),
	}
}
