// Package user is a service with all business logic.
package user

import (
	"context"
	"errors"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/hash"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/token"
	"os"
	"time"
)

var (
	ErrSignIn = errors.New("wrong user or password")
)

// ServiceAuth compose all functions.
type ServiceAuth struct {
	repo        repository.Auth
	tokenManger token.Manager
	hashManager *hash.Manager
}

// NewServiceAuth create service with all fields.
func NewServiceAuth(repo repository.Auth, tokenManger token.Manager, hashManager *hash.Manager) *ServiceAuth {
	return &ServiceAuth{
		repo:        repo,
		tokenManger: tokenManger,
		hashManager: hashManager,
	}
}

// CreateUser create new user in db and return error if any.
func (s *ServiceAuth) CreateUser(ctx context.Context, user domain.User) error {
	hashedPassword, err := s.hashManager.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.repo.CreateUser(ctx, user)
}

func (s *ServiceAuth) SignIn(ctx context.Context, user domain.User) (string, error) {
	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return "", err
	}

	userByEmail, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", ErrSignIn
	}

	if ok := s.hashManager.CheckPasswordHash(user.Password, userByEmail.Password); !ok {
		return "", ErrSignIn
	}

	return s.tokenManger.CreateToken(user.ID, duration)
}

// GetUser get user from db and return user and error if any.
func (s *ServiceAuth) GetUser(ctx context.Context, email string, password []byte) (domain.User, error) {
	return s.repo.GetUser(ctx, email, password)
}

// GetUserByEmail get user from db and return user and error if any.
func (s *ServiceAuth) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// DeleteUserById delete user from db by id and return error if any.
func (s *ServiceAuth) DeleteUserById(ctx context.Context, userID int) error {
	return s.repo.DeleteUserById(ctx, userID)
}

func (s *ServiceAuth) VerifyToken(ctx context.Context, token string) (*token.Payload, error) {
	return s.tokenManger.VerifyToken(token)
}
