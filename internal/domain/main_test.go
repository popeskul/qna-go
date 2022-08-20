package domain

import (
	"github.com/go-playground/validator/v10"
	"os"
	"testing"
)

var validate *validator.Validate

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
