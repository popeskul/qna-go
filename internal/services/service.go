// Package services implements all services.
package services

import (
	"context"
	"github.com/popeskul/cache"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services/tests"
	"github.com/popeskul/qna-go/internal/services/user"
	"github.com/popeskul/qna-go/internal/token"
)

// Auth interface is implemented by auth service.
type Auth interface {
	CreateUser(ctx context.Context, userInput domain.User) error
	GetUser(ctx context.Context, email string, password []byte) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

// Tests interface is implemented by tests service.
type Tests interface {
	CreateTest(ctx context.Context, userID int, test domain.Test) error
	GetTest(ctx context.Context, testID int) (domain.Test, error)
	GetAllTestsByUserID(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error)
	UpdateTestByID(ctx context.Context, testID int, test domain.Test) error
	DeleteTestByID(ctx context.Context, testID int) error
}

// Service struct is composed of all services.
type Service struct {
	Auth
	Tests
	TokenMaker token.Maker
	Cache      *cache.Cache
}

// NewService creates a new service with all services.
func NewService(repo *repository.Repository, tokenMaker token.Maker, cache *cache.Cache) *Service {
	return &Service{
		Auth:       user.NewServiceAuth(repo),
		Tests:      tests.NewServiceTests(repo, cache),
		TokenMaker: tokenMaker,
	}
}
