// Package user is a service with all business logic.
package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/hash"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/repository/sessions"
	"github.com/popeskul/qna-go/internal/token"
	"math/rand"
	"os"
	"time"
)

var (
	ErrSignIn = errors.New("wrong user or password")
)

// ServiceAuth compose all functions.
type ServiceAuth struct {
	repo           repository.Auth
	tokenManger    token.Manager
	hashManager    *hash.Manager
	sessionManager *sessions.RepositorySessions
}

// NewServiceAuth create service with all fields.
func NewServiceAuth(repo repository.Auth, tokenManger token.Manager, hashManager *hash.Manager, sessionManager *sessions.RepositorySessions) *ServiceAuth {
	return &ServiceAuth{
		repo:           repo,
		tokenManger:    tokenManger,
		hashManager:    hashManager,
		sessionManager: sessionManager,
	}
}

// CreateUser create new user in db and return error if any.
func (s *ServiceAuth) CreateUser(ctx context.Context, user domain.User) error {
	foundUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil && foundUser.Email == user.Email {
		return err
	}

	hashedPassword, err := s.hashManager.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.repo.CreateUser(ctx, user)
}

func (s *ServiceAuth) SignIn(ctx context.Context, user domain.User) (string, string, error) {
	userByEmail, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", "", ErrSignIn
	}

	if ok := s.hashManager.CheckPasswordHash(user.Password, userByEmail.Password); !ok {
		return "", "", ErrSignIn
	}

	return s.generateToken(ctx, userByEmail.ID)
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

func (s *ServiceAuth) GenerateAccessRefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.sessionManager.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token expired")
	}

	return s.generateToken(ctx, int(session.UserID))
}

func (s *ServiceAuth) generateToken(ctx context.Context, userID int) (string, string, error) {
	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.tokenManger.CreateToken(userID, duration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err = s.sessionManager.CreateRefreshToken(ctx, domain.RefreshSession{
		Token:     refreshToken,
		UserID:    int64(userID),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
