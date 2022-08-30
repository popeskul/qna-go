package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = fmt.Errorf("token is expired")
	ErrInvalidToken = fmt.Errorf("invalid token")
)

// Payload contains the payload data of the token.
type Payload struct {
	ID        uuid.UUID     `json:"id"`
	UserID    int           `json:"user_id"`
	IssuedAt  time.Duration `json:"issued_at"`
	ExpiredAt time.Duration `json:"expired_at"`
}

func NewPayload(userID int, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Duration(time.Now().Unix()),
		ExpiredAt: time.Duration(time.Now().Add(duration).Unix()),
	}, nil
}

// Valid check if the token payload is valid or not.
func (p *Payload) Valid() error {
	if time.Duration(time.Now().Unix()) > p.ExpiredAt {
		return ErrExpiredToken
	}
	return nil
}
