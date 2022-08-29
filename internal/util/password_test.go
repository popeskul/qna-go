package util

import "testing"

func TestHashPassword(t *testing.T) {
	password := RandomString(10)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	if hashedPassword == "" {
		t.Error("hashedPassword is empty")
	}

	match := CheckPassword(password, hashedPassword)
	if !match {
		t.Error("password does not match")
	}

	wrongPassword := RandomString(10)
	match = CheckPassword(wrongPassword, hashedPassword)
	if match {
		t.Error("password does not match")
	}
}
