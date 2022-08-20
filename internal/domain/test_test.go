package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/popeskul/qna-go/internal/util"
	"testing"
)

func Test_TestInput(t *testing.T) {
	validate = validator.New()

	tests := []struct {
		name    string
		fields  TestInput
		wantErr bool
	}{
		{
			name: "valid",
			fields: TestInput{
				Title: util.RandomString(10),
			},
			wantErr: false,
		},
		{
			name: "invalid title, with min length",
			fields: TestInput{
				Title: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
