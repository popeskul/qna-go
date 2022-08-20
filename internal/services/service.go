// Package services implements all services.
package services

import (
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services/auth"
	"github.com/popeskul/qna-go/internal/services/tests"
)

// Auth interface is implemented by auth service.
type Auth interface {
	CreateUser(userInput domain.SignUpInput) (int, error)
	GetUser(email, password string) (domain.User, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Tests interface is implemented by tests service.
type Tests interface {
	CreateTest(userID int, testInput domain.TestInput) (int, error)
	GetTest(testID int) (domain.Test, error)
	UpdateTestByID(testID int, testInput domain.TestInput) error
	DeleteTestByID(testID int) error
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
