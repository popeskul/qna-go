package services

import (
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services/auth"
	"github.com/popeskul/qna-go/internal/services/tests"
)

type Auth interface {
	CreateUser(userInput domain.SignUpInput) (int, error)
	GetUser(email, password string) (domain.User, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Tests interface {
	CreateTest(userID int, testInput domain.TestInput) (int, error)
	GetTest(testID int) (domain.Test, error)
	UpdateTestByID(testID int, testInput domain.TestInput) error
	DeleteTestByID(testID int) error
}

type Service struct {
	Auth
	Tests
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:  auth.NewServiceAuth(repo),
		Tests: tests.NewServiceTests(repo),
	}
}
