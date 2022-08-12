package domain

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

var validate *validator.Validate

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
				Name:              "John Doe",
				Email:             "test@mail.com",
				EncryptedPassword: "password",
			},
			wantErr: false,
		},
		{
			name: "invalid name, with min length",
			fields: SignUpInput{
				Name:              "",
				Email:             "test@mail.com",
				EncryptedPassword: "password",
			},
			wantErr: true,
		},
		{
			name: "invalid name with min length",
			fields: SignUpInput{
				Name:              "",
				Email:             "test@mail.com",
				EncryptedPassword: "password",
			},
			wantErr: true,
		},
		{
			name: "invalid name with max length",
			fields: SignUpInput{
				Name:              "of15XzFMJ0CDkVBGf5bhjqXM5SgN3IVLE3L5f2n2t6S13w2WmaGER1d5brKxLLWiODCbpvQKZzmX8L5vHAiZ7KXnuJnNB5BT74irP1yoKJ2JKWDK2l1NgHAa63Ewu0OWg86GuFoNql6pRradtW10AOUsDSFwE8rqLIo3GWjy3UNPxCI606A52gF1pKUQRnnWtCMbwKnufvs2LdijZbkFNuurtY3jTQ3CHHjaph5GBAVffJhJw8RMTVI3NnywnzEz",
				Email:             "test@mail.com",
				EncryptedPassword: "password",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: SignUpInput{
				Name:              "12345",
				Email:             "",
				EncryptedPassword: "password",
			},
			wantErr: true,
		},
		{
			name: "invalid password with min length",
			fields: SignUpInput{
				Name:              "12345",
				Email:             "test@mail.com",
				EncryptedPassword: "",
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
