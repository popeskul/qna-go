package domain

import (
	"github.com/go-playground/validator/v10"
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
				Title:    "test",
				AuthorID: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid title, with min length",
			fields: TestInput{
				Title:    "",
				AuthorID: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid without author id",
			fields: TestInput{
				Title: "123123123",
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
