package auth

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"time"
)

const (
	salt       = "asd13ewd"
	signingKey = "asd13ewd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type ServiceAuth struct {
	repo repository.Auth
}

func NewServiceAuth(repo repository.Auth) *ServiceAuth {
	return &ServiceAuth{
		repo: repo,
	}
}

func (s *ServiceAuth) CreateUser(input domain.SignUpInput) (int, error) {
	passwordPasha, err := generatePasswordHash(input.Password)
	if err != nil {
		return 0, err
	}
	
	input.Password = passwordPasha
	return s.repo.CreateUser(input)
}

func (s *ServiceAuth) GetUser(email, password string) (domain.User, error) {
	return s.repo.GetUser(email, password)
}

func (s *ServiceAuth) GenerateToken(username, password string) (string, error) {
	passwordHash, err := generatePasswordHash(password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetUser(username, passwordHash)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *ServiceAuth) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserId, nil
}

func (s *ServiceAuth) DeleteUserById(userID int) error {
	return s.repo.DeleteUserById(userID)
}

func generatePasswordHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil
}
