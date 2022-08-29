package token

import (
	"github.com/popeskul/qna-go/internal/util"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	userID := 1
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := jwtMaker.CreateToken(userID, duration)
	if err != nil {
		t.Fatal(err)
	}

	payload, err := jwtMaker.VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if payload.UserID != userID {
		t.Fatal("user_id is not correct")
	}
	if payload.IssuedAt < time.Duration(issuedAt.Unix()) {
		t.Fatal("issued_at is not correct")
	}
	if payload.ExpiredAt > time.Duration(expiredAt.Unix()) {
		t.Fatal("expired_at is not correct")
	}
}

func TestExpiredJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatal(err)
	}

	token, err := jwtMaker.CreateToken(1, -time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	payload, err := jwtMaker.VerifyToken(token)
	if err != ErrExpiredToken {
		t.Fatal("token is not expired")
	}
	if payload != nil {
		t.Fatal("payload is not nil")
	}
}
