package tests

import (
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/repository/tests"
	"github.com/popeskul/qna-go/internal/util"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

var (
	mockDB   *sql.DB
	mockRepo *repository.Repository
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

	os.Exit(m.Run())
}

func TestServiceTests_CreateTest(t *testing.T) {
	ctx := context.Background()
	u := randomTest()
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
				input:  u,
				userID: mockUserID,
			},
			want: want{
				title: u.Title,
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
				err: tests.ErrEmptyTitle,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			testID, err := mockRepo.CreateTest(ctx, tt.args.userID, tt.args.input)

			if err != tt.want.err {
				t.Errorf("ServiceTests.CreateTest() error = %v, wantErr %v", err, tt.want.err)
				return
			}

			t.Cleanup(func() {
				helperDeleteTestByID(t, testID)
			})
		})
	}
}

func TestServiceTests_UpdateTestById(t *testing.T) {
	ctx := context.Background()
	mockUserID := 1
	testID, err := mockRepo.CreateTest(ctx, mockUserID, randomTest())
	if err != nil {
		t.Errorf("error creating test: %v", err)
	}

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
			name: "Success: Update test",
			args: args{
				repo: mockRepo,
				input: domain.TestInput{
					Title: "Test title updated",
				},
				userID: mockUserID,
			},
			want: want{
				title: "Test title updated",
			},
		},
		{
			name: "Fail: Update test with empty title",
			args: args{
				repo: mockRepo,
				input: domain.TestInput{
					Title: "",
				},
				userID: mockUserID,
			},
			want: want{
				err: tests.ErrEmptyTitle,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err = mockRepo.UpdateTestById(ctx, tt.args.userID, tt.args.input)

			if err != tt.want.err {
				t.Errorf("ServiceTests.UpdateTestById() error = %v, wantErr %v", err, tt.want.err)
				return
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteTestByID(t, testID)
	})
}

func TestServiceTests_DeleteTestById(t *testing.T) {
	ctx := context.Background()
	mockUserID := 1

	type args struct {
		repo   *repository.Repository
		userID int
	}
	type want struct {
		err error
	}
	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: Delete test",
			args: args{
				repo:   mockRepo,
				userID: mockUserID,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Fail: Delete test",
			args: args{
				repo:   mockRepo,
				userID: 0,
			},
			want: want{
				err: tests.ErrTestAuthorIDEmpty,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			testID, err := mockRepo.CreateTest(ctx, tt.args.userID, randomTest())
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("ServiceTests.DeleteTestById() error = %v, wantErr %v", err, tt.want.err)
			}

			err = mockRepo.DeleteTestById(ctx, tt.args.userID)
			if err != tt.want.err {
				t.Errorf("ServiceTests.DeleteTestById() error = %v, wantErr %v", err, tt.want.err)
			}

			t.Cleanup(func() {
				helperDeleteTestByID(t, testID)
			})
		})
	}
}

func randomTest() domain.TestInput {
	return domain.TestInput{
		Title: util.RandomString(10),
	}
}

func helperDeleteTestByID(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE id = $1", id); err != nil {
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
