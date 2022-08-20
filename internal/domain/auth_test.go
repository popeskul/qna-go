package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/popeskul/qna-go/internal/util"
	"testing"
)

func Test_SignUpInput(t *testing.T) {
	validate = validator.New()

	tests := []struct {
		name    string
		fields  SignUpInput
		wantErr bool
	}{
		{
			name: "valid",
			fields: SignUpInput{
				Name:     util.RandomString(10),
				Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				Password: util.RandomString(10),
			},
			wantErr: false,
		},
		{
			name: "invalid name, with min length",
			fields: SignUpInput{
				Name:     "",
				Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				Password: util.RandomString(10),
			},
			wantErr: true,
		},
		{
			name: "invalid name with min length",
			fields: SignUpInput{
				Name:     "",
				Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				Password: util.RandomString(10),
			},
			wantErr: true,
		},
		{
			name: "invalid name with max length",
			fields: SignUpInput{
				Name:     util.RandomString(256),
				Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				Password: util.RandomString(10),
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: SignUpInput{
				Name:     util.RandomString(10),
				Email:    "",
				Password: util.RandomString(10),
			},
			wantErr: true,
		},
		{
			name: "invalid password with min length",
			fields: SignUpInput{
				Name:     util.RandomString(10),
				Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
