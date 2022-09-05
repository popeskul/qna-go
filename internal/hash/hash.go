package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHash = fmt.Errorf("error hashing")
)

type Manager struct {
	salt string
}

func NewHash(salt string) (*Manager, error) {
	if salt == "" {
		return &Manager{}, ErrHash
	}

	return &Manager{
		salt: salt,
	}, nil
}

func (h *Manager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHash
	}

	return string(hashedPassword), nil
}

func (h *Manager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
