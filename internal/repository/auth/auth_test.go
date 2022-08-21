package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var mockDB *sql.DB
var mockRepo *RepositoryAuth

func TestMain(m *testing.M) {
	if err := util.ChangeDir("./../../../"); err != nil {
		log.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := newDBConnection(cfg)
	mockDB = db
	if err != nil {
		log.Fatal(err)
	}

	mockRepo = NewRepoAuth(mockDB)

	os.Exit(m.Run())
}

func TestRepositoryAuth_CreateUser(t *testing.T) {
	ctx := context.Background()

	u := randomUser()
	userID, err := mockRepo.CreateUser(ctx, u)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		repo *RepositoryAuth
	}
	type args struct {
		u domain.SignUpInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "CreateUser",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				u: randomUser(),
			},
			wantErr: false,
		},
		{
			name: "CreateUser_DuplicateEmail",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				u: domain.SignUpInput{
					Name:  util.RandomString(10),
					Email: u.Email,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := tt.fields.repo.CreateUser(ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryAuth.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Cleanup(func() {
				helperDeleteUserByID(t, id)
			})
		})
	}

	t.Cleanup(func() {
		if err = mockRepo.DeleteUserById(ctx, userID); err != nil {
			t.Error(err)
		}
	})
}

func TestRepositoryAuth_GetUser(t *testing.T) {
	ctx := context.Background()
	u := randomUser()
	userId, err := mockRepo.CreateUser(ctx, u)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		repo *RepositoryAuth
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "GetUser",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				email:    u.Email,
				password: u.Password,
			},
			want: userId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.repo.GetUser(ctx, tt.args.email, tt.args.password)
			if err != nil {
				t.Error(err)
			}

			if got.ID != tt.want {
				t.Errorf("RepositoryAuth.GetUser() = %v, want %v", got.ID, tt.want)
			}
		})
	}

	t.Cleanup(func() {
		if err = mockRepo.DeleteUserById(ctx, userId); err != nil {
			t.Error(err)
		}
	})
}

func TestRepositoryAuth_DeleteUserById(t *testing.T) {
	ctx := context.Background()
	user := domain.SignUpInput{
		Name:     "John Doe",
		Email:    "TestRepositoryAuth_DeleteUserById@mail.com",
		Password: "12345",
	}
	userId, err := mockRepo.CreateUser(ctx, user)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		repo *RepositoryAuth
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: DeleteUserById",
			args: args{
				repo: mockRepo,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Fail: DeleteUserById",
			args: args{
				repo: mockRepo,
			},
			want: want{
				err: errors.New("record not found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.repo.DeleteUserById(ctx, userId); err != nil {
				t.Error(err)
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userId)
	})
}

func randomUser() domain.SignUpInput {
	return domain.SignUpInput{
		Name:     util.RandomString(10),
		Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
		Password: util.RandomString(10),
	}
}

func helperDeleteUserByID(t *testing.T, id int) {
	ctx := context.Background()
	if err := mockRepo.DeleteUserById(ctx, id); err != nil {
		t.Error(err)
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg, err := config.New("configs", "test.config")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return cfg, nil
}
