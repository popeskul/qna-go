package tests

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/repository/tests"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

var (
	mockDB   *sql.DB
	mockRepo *repository.Repository
)

func TestMain(m *testing.M) {
	if err := changeDirToRoot(); err != nil {
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

	os.Exit(m.Run())
}

func TestServiceTests_CreateTest(t *testing.T) {
	mockTest := domain.TestInput{
		Title: "Test title",
	}
	mockUserID := 1

	type args struct {
		repo   *repository.Repository
		input  domain.TestInput
		userID int
	}
	type want struct {
		title string
		err   error
	}

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: Create test",
			args: args{
				repo:   mockRepo,
				input:  mockTest,
				userID: mockUserID,
			},
			want: want{
				title: mockTest.Title,
			},
		},
		{
			name: "Fail: Create test",
			args: args{
				repo: mockRepo,
				input: domain.TestInput{
					Title: "",
				},
				userID: mockUserID,
			},
			want: want{
				err: tests.ErrTestTitleEmpty,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockRepo.CreateTest(tt.args.userID, tt.args.input)

			if err != tt.want.err {
				t.Errorf("ServiceTests.CreateTest() error = %v, wantErr %v", err, tt.want.err)
				return
			}

			t.Cleanup(func() {
				helperDeleteTestByTitle(t, tt.args.input.Title)
			})
		})
	}
}

func helperDeleteTestByTitle(t *testing.T, title string) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE title = $1", title); err != nil {
		t.Errorf("error deleting test: %v", err)
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

func changeDirToRoot() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "./../../../")
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}
