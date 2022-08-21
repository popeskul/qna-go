package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/util"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	mockDB      *sql.DB
	mockRepo    *repository.Repository
	mockService *ServiceAuth
)

func TestMain(m *testing.M) {
	if err := util.ChangeDir("./../../../"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	db, err := newDBConnection(cfg)
	mockDB = db
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer mockDB.Close()

	mockRepo = repository.NewRepository(mockDB)
	mockService = NewServiceAuth(mockRepo)

	os.Exit(m.Run())
}

func TestServiceAuth_CreateUser(t *testing.T) {
	ctx := context.Background()
	u := randomUser()

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
				repo: mockRepo,
			},
			args: args{
				input: u,
			},
			err: nil,
		},
		{
			name: "CreateUser_DuplicateEmail",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				input: u,
			},
			err: errors.New("duplicate key value violates unique constraint"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := NewServiceAuth(tt.fields.repo).CreateUser(ctx, tt.args.input)
			if err != nil {
				if !strings.Contains(tt.err.Error(), err.Error()) {
					t.Errorf("ServiceAuth.CreateUser() error = %v, wantErr %v", err, tt.err)
				}
			}

			t.Cleanup(func() {
				helperDeleteUserByID(t, userID)
			})
		})
	}
}

func TestServiceAuth_GetUser(t *testing.T) {
	ctx := context.Background()
	u := randomUser()
	userID, err := mockRepo.CreateUser(ctx, u)
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
				repo: mockRepo,
			},
			args: args{
				email:    u.Email,
				password: u.Password,
			},
			want: &domain.User{
				Email:             u.Email,
				EncryptedPassword: u.Password,
			},
		},
		{
			name: "fail",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				password: util.RandomString(10),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceAuth(tt.fields.repo)
			got, err := s.GetUser(ctx, tt.args.email, tt.args.password)

			if err != nil {
				if tt.want != nil && tt.want.Email != got.Email {
					t.Errorf("ServiceAuth.GetUser() error = %v, wantErr %v", err, tt.want)
				}
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
	})
}

func TestServiceAuth_GenerateToken(t *testing.T) {
	ctx := context.Background()
	u := randomUser()
	userID, err := mockService.CreateUser(ctx, u)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		user domain.SignInInput
	}
	tests := []struct {
		name      string
		args      args
		wantError error
	}{
		{
			name: "success",
			args: args{
				user: domain.SignInInput{
					Email:    u.Email,
					Password: u.Password,
				},
			},
			wantError: nil,
		},
		{
			name: "fail",
			args: args{
				user: domain.SignInInput{
					Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
					Password: util.RandomString(10),
				},
			},
			wantError: errors.New("sql: no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = mockService.GenerateToken(ctx, tt.args.user.Email, tt.args.user.Password)

			if err != nil && tt.wantError != nil {
				if !strings.Contains(err.Error(), tt.wantError.Error()) {
					t.Errorf("ServiceAuth.GenerateToken() error = %v, wantErr = %v", err, tt.wantError)
				}
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
	})
}

func TestServiceAuth_generatePassword(t *testing.T) {
	token, err := generatePasswordHash(util.RandomString(10))
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("token is empty")
	}
}

func helperDeleteUserByID(t *testing.T, userID int) {
	t.Helper()

	_, err := mockDB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Error(err)
	}
}

func randomUser() domain.SignUpInput {
	return domain.SignUpInput{
		Name:     util.RandomString(10),
		Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
		Password: util.RandomString(10),
	}
}

func newDBConnection(cfg *config.Config) (*sql.DB, error) {
	return postgres.NewPostgresConnection(db.ConfigDB{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
}

func loadConfig() (*config.Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg, err := config.New("configs", "test.config")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return cfg, nil
}
