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
	"github.com/popeskul/qna-go/internal/hash"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/token"
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
	if err := util.ChangeDir("../../"); err != nil {
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

	pasetoMaker, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		log.Fatal(err)
	}

	hashManager, err := hash.NewHash(cfg.HashSalt)
	if err != nil {
		log.Fatal(err)
	}

	mockRepo = repository.NewRepository(mockDB)
	mockService = NewServiceAuth(mockRepo, pasetoMaker, hashManager)

	os.Exit(m.Run())
}

func TestServiceAuth_CreateUser(t *testing.T) {
	ctx := context.Background()
	u := randomUser()

	type fields struct {
		repo *repository.Repository
	}
	type args struct {
		input domain.User
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
			err := mockService.CreateUser(ctx, tt.args.input)
			if err != nil {
				if !strings.Contains(tt.err.Error(), err.Error()) {
					t.Errorf("ServiceAuth.CreateUser() error = %v, wantErr %v", err, tt.err)
				}
			}

			t.Cleanup(func() {
				if err == nil {
					userID, err := findUserIDByEmail(u.Email)
					if err != nil {
						t.Fatalf("Some error occured. Err: %s", err)
					}
					helperDeleteUserByID(t, userID)
				}
			})
		})
	}
}

func TestServiceAuth_GetUser(t *testing.T) {
	ctx := context.Background()
	u := randomUser()

	if err := mockRepo.CreateUser(ctx, u); err != nil {
		t.Fatalf("Some error occured. Err: %s", err)
	}

	userID, err := findUserIDByEmail(u.Email)
	if err != nil {
		t.Fatalf("Some error occured. Err: %s", err)
	}

	type fields struct {
		repo *repository.Repository
	}
	type args struct {
		email    string
		password []byte
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
				password: []byte(u.Password),
			},
			want: &domain.User{
				Email:    u.Email,
				Password: u.Password,
			},
		},
		{
			name: "fail",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
				password: []byte(util.RandomString(10)),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockService.GetUser(ctx, tt.args.email, tt.args.password)

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

func helperDeleteUserByID(t *testing.T, userID int) {
	t.Helper()

	_, err := mockDB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Fatalf("Some error occured. Err: %s", err)
	}
}

func findUserIDByEmail(email string) (int, error) {
	var userID int
	err := mockDB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func randomUser() domain.User {
	return domain.User{
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
	cfg.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	cfg.HashSalt = os.Getenv("HASH_SALT")

	return cfg, nil
}
