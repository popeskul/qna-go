// Package repository is a struct that contains the repository.
// This place define interface for the repository: Auth, Tests.
package repository

import (
	"context"
	"database/sql"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/auth"
	"github.com/popeskul/qna-go/internal/repository/tests"
)

// Auth interface is implemented by the auth repository.
type Auth interface {
	CreateUser(ctx context.Context, userInput domain.SignUpInput) (int, error)
	GetUser(ctx context.Context, email, password string) (domain.User, error)
	DeleteUserById(ctx context.Context, userID int) error
}

// Tests interface is implemented by the test repository.
type Tests interface {
	CreateTest(ctx context.Context, userID int, testInput domain.TestInput) (int, error)
	GetTest(ctx context.Context, testID int) (domain.Test, error)
	GetAllTestsByCurrentUser(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error)
	UpdateTestById(ctx context.Context, testID int, testInput domain.TestInput) error
	DeleteTestById(ctx context.Context, testID int) error
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
