// Package repository is a struct that contains the repository.
// This place define interface for the repository: Auth, Tests.
package repository

import (
	"context"
	"database/sql"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/tests"
	"github.com/popeskul/qna-go/internal/repository/user"
)

// Auth interface is implemented by the auth repository.
type Auth interface {
	CreateUser(ctx context.Context, userInput domain.User) error
	GetUser(ctx context.Context, email string, password []byte) (domain.User, error)
	DeleteUserById(ctx context.Context, userID int) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

// Tests interface is implemented by the test repository.
type Tests interface {
	CreateTest(ctx context.Context, userID int, test domain.Test) error
	GetTest(ctx context.Context, testID int) (domain.Test, error)
	GetAllTestsByUserID(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error)
	UpdateTestById(ctx context.Context, testID int, test domain.Test) error
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
		Auth:  user.NewRepoAuth(db),
		Tests: tests.NewRepoTests(db),
	}
}
