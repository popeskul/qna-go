package auth

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

var dbConn *sql.DB

func TestMain(m *testing.M) {
	if err := changeDirToRoot(); err != nil {
		log.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := newDBConnection(cfg)
	dbConn = db
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestRepositoryAuth_CreateUser(t *testing.T) {
	repo := NewRepoAuth(dbConn)

	mockSimpleEmail := "testting1@test.com"
	mockUniqueEmail := "test_unique_email@mail.com"
	mockPassword := "12345"

	// Create user with simple email
	_, err := repo.GetUser(mockUniqueEmail, mockPassword)
	if err != nil {
		// create single user for testing duplicate email
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
				repo: repo,
			},
			args: args{
				u: domain.SignUpInput{
					Name:              "John Doe",
					Email:             mockSimpleEmail,
					EncryptedPassword: "encrypted_password",
				},
			},
			wantErr: false,
		},
		{
			name: "CreateUser_DuplicateEmail",
			fields: fields{
				repo: repo,
			},
			args: args{
				u: domain.SignUpInput{
					Name:  "John Doe",
					Email: mockUniqueEmail,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fields.repo.CreateUser(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryAuth.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Cleanup
			cleanupQuery := fmt.Sprintf("DELETE FROM users WHERE email = $1")
			if _, err = dbConn.Exec(cleanupQuery, tt.args.u.Email); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRepositoryAuth_GetUser(t *testing.T) {
	repo := NewRepoAuth(dbConn)

	mockEmail := "testting2@test.com"
	mockPassword := "12345"
	userId, err := repo.CreateUser(domain.SignUpInput{
		Name:              "John Doe",
		Email:             mockEmail,
		EncryptedPassword: mockPassword,
	})
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
				repo: repo,
			},
			args: args{
				email:    mockEmail,
				password: mockPassword,
			},
			want: userId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.repo.GetUser(tt.args.email, tt.args.password)
			if err != nil {
				t.Error(err)
			}

			if got.ID != tt.want {
				t.Errorf("RepositoryAuth.GetUser() = %v, want %v", got.ID, tt.want)
			}

			// Cleanup
			cleanupQuery := fmt.Sprintf("DELETE FROM users WHERE email = $1")
			if _, err = dbConn.Exec(cleanupQuery, tt.args.email); err != nil {
				t.Error(err)
			}
		})
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

	cfg, err := config.New("configs", "config")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return cfg, nil
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
