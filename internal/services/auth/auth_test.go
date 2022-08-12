package auth

import (
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestServiceAuth_CreateUser(t *testing.T) {
	mockEmail := "TestServiceAuth_CreateUser@example.com"
	mockUniqueEmail := "test_unique_email@mail.com"
	mockPassword := "12345"

	db, err := newDBConnection()
	if err != nil {
		t.Error(err)
	}

	repo := repository.NewRepository(db)

	// Create seed user for testing duplicate email
	_, err = repo.GetUser(mockUniqueEmail, mockPassword)
	if err != nil {
		_, err = repo.CreateUser(domain.SignUpInput{
			Email:             mockUniqueEmail,
			EncryptedPassword: mockPassword,
			Name:              "test",
		})
		if err != nil {
			t.Error(err)
		}
	}

	type fields struct {
		repo *repository.Repository
	}
	type args struct {
		input domain.SignUpInput
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		err    error
	}{
		{
			name: "CreateUser",
			fields: fields{
				repo: repo,
			},
			args: args{
				input: domain.SignUpInput{
					Email:             mockEmail,
					EncryptedPassword: mockPassword,
				},
			},
			err: nil,
		},
		{
			name: "CreateUser_DuplicateEmail",
			fields: fields{
				repo: repo,
			},
			args: args{
				input: domain.SignUpInput{
					Email:             mockUniqueEmail,
					EncryptedPassword: mockPassword,
				},
			},
			err: errors.New("duplicate key value violates unique constraint"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceAuth(tt.fields.repo)
			_, err := s.CreateUser(tt.args.input)
			if err != nil {
				if strings.Contains(tt.err.Error(), err.Error()) {
					t.Errorf("ServiceAuth.CreateUser() error = %v, wantErr %v", err, tt.err)
				}
			}

			// cleanup
			_, err = db.Exec("DELETE FROM users WHERE email = $1", tt.args.input.Email)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestServiceAuth_GetUser(t *testing.T) {
	mockEmail := "TestServiceAuth_GetUser@test.com"
	mockPassword := "12345"

	db, err := newDBConnection()
	if err != nil {
		t.Error(err)
	}

	repo := repository.NewRepository(db)
	_, err = repo.CreateUser(domain.SignUpInput{
		Email:             mockEmail,
		EncryptedPassword: mockPassword,
	})
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		repo *repository.Repository
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.User
	}{
		{
			name: "success",
			fields: fields{
				repo: repo,
			},
			args: args{
				email:    mockEmail,
				password: mockPassword,
			},
			want: &domain.User{
				Email:             mockEmail,
				EncryptedPassword: mockPassword,
			},
		},
		{
			name: "fail",
			fields: fields{
				repo: repo,
			},
			args: args{
				email:    "bad@mail.com",
				password: "123123123",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceAuth(tt.fields.repo)
			got, err := s.GetUser(tt.args.email, tt.args.password)

			if err != nil {
				if tt.want != nil && tt.want.Email != got.Email {
					t.Errorf("ServiceAuth.GetUser() error = %v, wantErr %v", err, tt.want)
				}
			}

			// cleanup
			_, err = db.Exec("DELETE FROM users WHERE email = $1", tt.args.email)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestServiceAuth_GenerateToken(t *testing.T) {
	mockUser := domain.User{
		Email:             "TestServiceAuth_GenerateToken@mail.com",
		EncryptedPassword: "123456",
	}

	// Create seed user for testing duplicate email
	db, err := newDBConnection()
	if err != nil {
		t.Error(err)
	}

	repo := repository.NewRepository(db)
	service := NewServiceAuth(repo)

	_, _ = service.CreateUser(domain.SignUpInput{
		Email:             mockUser.Email,
		EncryptedPassword: mockUser.EncryptedPassword,
	})

	type fields struct {
		repo *repository.Repository
	}
	type args struct {
		user *domain.User
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError error
	}{
		{
			name: "success",
			fields: fields{
				repo: repository.NewRepository(db),
			},
			args: args{
				user: &mockUser,
			},
			wantError: nil,
		},
		{
			name: "fail",
			fields: fields{
				repo: repository.NewRepository(db),
			},
			args: args{
				user: &domain.User{
					Email:             "wrong@mail.com",
					EncryptedPassword: "asd123",
				},
			},
			wantError: errors.New("sql: no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceAuth(tt.fields.repo)
			_, err := s.GenerateToken(tt.args.user.Email, tt.args.user.EncryptedPassword)

			if err != nil && tt.wantError != nil {
				if !strings.Contains(err.Error(), tt.wantError.Error()) {
					t.Errorf("ServiceAuth.GenerateToken() error = %v, wantErr = %v", err, tt.wantError)
				}
			}

			// cleanup
			t.Cleanup(func() {
				_, err := db.Exec("DELETE FROM users WHERE email = $1", tt.args.user.Email)
				if err != nil {
					t.Error(err)
				}
			})
		})
	}
}

func TestServiceAuth_generatePassword(t *testing.T) {
	token := generatePasswordHash("123456")
	if token == "" {
		t.Error("token is empty")
	}
}

func newDBConnection() (*sql.DB, error) {
	err := changeDirToRoot()
	if err != nil {
		return nil, err
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg, err := config.New("configs", "config")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return postgres.NewPostgresConnection(db.ConfigDB{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
}

func changeDirToRoot() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "./../../../")
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}
