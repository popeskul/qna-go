// Package services implements all services.
package services

import (
	"context"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/hash"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/repository/sessions"
	"github.com/popeskul/qna-go/internal/services/tests"
	"github.com/popeskul/qna-go/internal/services/user"
	"github.com/popeskul/qna-go/internal/token"
)

// Auth interface is implemented by auth service.
type Auth interface {
	CreateUser(ctx context.Context, userInput domain.User) error
	SignIn(ctx context.Context, userInput domain.User) (string, string, error)
	GetUser(ctx context.Context, email string, password []byte) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	VerifyToken(ctx context.Context, token string) (*token.Payload, error)
	GenerateAccessRefreshTokens(ctx context.Context, token string) (string, string, error)
}

// Sessions interface is implemented by sessions' repository.
type Sessions interface {
	CreateRefreshToken(ctx context.Context, token domain.RefreshSession) error
	GetRefreshToken(ctx context.Context, token string) (domain.RefreshSession, error)
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
	Sessions
}

// NewService creates a new service with all services.
func NewService(repo *repository.Repository, tokenManager token.Manager, hashManager *hash.Manager, sessionManager *sessions.RepositorySessions) *Service {
	return &Service{
		Auth:  user.NewServiceAuth(repo, tokenManager, hashManager, sessionManager),
		Tests: tests.NewServiceTests(repo),
	}
}
