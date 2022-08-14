package domain

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

var validate *validator.Validate

func TestMain(m *testing.M) {
	m.Run()
}
