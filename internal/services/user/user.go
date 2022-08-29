// Package user is a service with all business logic.
package user

import (
	"context"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
)

// ServiceAuth compose all functions.
type ServiceAuth struct {
	repo repository.Auth
}

// NewServiceAuth create service with all fields.
func NewServiceAuth(repo repository.Auth) *ServiceAuth {
	return &ServiceAuth{
		repo: repo,
	}
}

// CreateUser create new user in db and return error if any.
func (s *ServiceAuth) CreateUser(ctx context.Context, user domain.User) error {
	return s.repo.CreateUser(ctx, user)
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
