package hash

import (
	"github.com/popeskul/qna-go/internal/util"
	"testing"
)

var (
	salt     = "asdadfrg23rkl"
	password = "12345"
)

func TestManager_HashPassword(t *testing.T) {
	manager, err := NewHash(salt)
	if err != nil {
		t.Fatal(err)
	}

	hashedPassword, err := manager.HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	if hashedPassword == "" {
		t.Error("hashedPassword is empty")
	}

	match := manager.CheckPasswordHash(password, hashedPassword)
	if !match {
		t.Error("password does not match")
	}

	wrongPassword := util.RandomString(10)
	match = manager.CheckPasswordHash(wrongPassword, hashedPassword)
	if match {
		t.Error("password does not match")
	}
}
