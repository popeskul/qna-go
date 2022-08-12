package services

import (
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services/auth"
)

type Auth interface {
	CreateUser(userInput domain.SignUpInput) (int, error)
	GetUser(email, password string) (domain.User, error)
	GenerateToken(username, password string) (string, error)
}

type Service struct {
	Auth
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth: auth.NewServiceAuth(repo),
	}
}
