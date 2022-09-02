package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretLength = 32

type JWTMaker struct {
	secretKey string
}

// NewJWTMaker returns a new JWTMaker.
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretLength {
		return nil, ErrSecretIsTooShort
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

// CreateToken returns a new token for specific userID and duration.
func (maker *JWTMaker) CreateToken(userID int, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken returns the payload of the token if the token is valid.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	claims, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := claims.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
