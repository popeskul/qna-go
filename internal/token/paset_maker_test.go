package token

import (
	"github.com/popeskul/qna-go/internal/util"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	userID := 1
	pasetoMaker, err := NewPasetoManager(util.RandomString(32))
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := pasetoMaker.CreateToken(userID, duration)
	if err != nil {
		t.Fatal(err)
	}

	payload, err := pasetoMaker.VerifyToken(token)
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

func TestExpiredPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoManager(util.RandomString(32))
	if err != nil {
		t.Fatal(err)
	}

	wrongDuration := -time.Minute
	token, err := pasetoMaker.CreateToken(1, wrongDuration)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	payload, err := pasetoMaker.VerifyToken(token)
	if err != ErrExpiredToken {
		t.Fatal("token is not expired")
	}
	if payload != nil {
		t.Fatal("payload is not nil")
	}
}
