// Package services implements all services.
package services

import (
	"context"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services/auth"
	"github.com/popeskul/qna-go/internal/services/tests"
)

// Auth interface is implemented by auth service.
type Auth interface {
	CreateUser(ctx context.Context, userInput domain.SignUpInput) error
	GetUser(ctx context.Context, email, password string) (domain.User, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Tests interface is implemented by tests service.
type Tests interface {
	CreateTest(ctx context.Context, userID int, testInput domain.TestInput) error
	GetTest(ctx context.Context, testID int) (domain.Test, error)
	GetAllTestsByCurrentUser(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error)
	UpdateTestByID(ctx context.Context, testID int, testInput domain.TestInput) error
	DeleteTestByID(ctx context.Context, testID int) error
}

// Service struct is composed of all services.
type Service struct {
	Auth
	Tests
}

// NewService creates a new service with all services.
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:  auth.NewServiceAuth(repo),
		Tests: tests.NewServiceTests(repo),
	}
}
