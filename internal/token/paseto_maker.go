package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

// PasetoMaker is a paseto maker.
type PasetoMaker struct {
	paseto        *paseto.V2
	symmetrickKey []byte
}

// NewPasetoMaker create new paseto maker with symmetric key and return paseto maker and error if any.
func NewPasetoMaker(symmetrickKey string) (Maker, error) {
	if len(symmetrickKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("symmetric key is too short")
	}

	return &PasetoMaker{
		paseto:        paseto.NewV2(),
		symmetrickKey: []byte(symmetrickKey),
	}, nil
}

// CreateToken create new token and return token and error if any.
func (maker *PasetoMaker) CreateToken(userID int, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetrickKey, payload, nil)
}

// VerifyToken verify token and return payload and error if any.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetrickKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
