package tests

import (
	"context"
	"database/sql"
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
var mockRepo *RepositoryTests

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

	mockRepo = NewRepoTests(mockDB)

	os.Exit(m.Run())
}

func TestRepositoryTests_CreateTest(t *testing.T) {
	ctx := context.Background()
	mockUserID := 1
	testID := helperCreateTest(t, randomTest(), 1)

	type args struct {
		repo   *RepositoryTests
		input  domain.TestInput
		userID int
	}
	type want struct {
		id  int
		err error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: create test",
			args: args{
				repo:   mockRepo,
				input:  randomTest(),
				userID: mockUserID,
			},
			want: want{
				id:  testID + 1,
				err: nil,
			},
		},
		{
			name: "Fail: create test",
			args: args{
				repo: mockRepo,
				input: domain.TestInput{
					Title: "",
				},
			},
			want: want{
				err: ErrDuplicateAuthorID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := tt.args.repo.CreateTest(ctx, tt.args.userID, tt.args.input)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("RepositoryTests.CreateTest() error = %v, wantErr %v", err, tt.want.err)
				return
			}

			t.Cleanup(func() {
				helperDeleteTest(t, id)
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteTest(t, testID)
	})
}

func TestRepositoryTests_UpdateTestById(t *testing.T) {
	ctx := context.Background()
	mockTestAuthorID := 1
	createdID := helperCreateTest(t, randomTest(), mockTestAuthorID)

	type args struct {
		repo   *RepositoryTests
		testID int
		input  domain.TestInput
	}
	type want struct {
		rest domain.Test
		err  error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: update test",
			args: args{
				repo:   mockRepo,
				testID: createdID,
				input: domain.TestInput{
					Title: "Test title updated",
				},
			},
			want: want{
				rest: domain.Test{
					Title: "Test title updated",
				},
				err: nil,
			},
		},
		{
			name: "Fail: update test with empty title",
			args: args{
				repo:   mockRepo,
				testID: createdID,
				input: domain.TestInput{
					Title: "",
				},
			},
			want: want{
				err: ErrEmptyTitle,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.repo.UpdateTestById(ctx, tt.args.testID, tt.args.input)
			if err != tt.want.err {
				t.Errorf("RepositoryTests.UpdateTestById() error = %v, wantErr %v", err, tt.want.err)
			}

			test, _ := tt.args.repo.GetTest(ctx, createdID)
			if test.Title != tt.want.rest.Title && err == nil {
				t.Errorf("RepositoryTests.UpdateTestById() error = %v, wantErr %v", test.Title, tt.want.rest.Title)
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteTest(t, createdID)
	})
}

func TestRepositoryTests_DeleteTestById(t *testing.T) {
	ctx := context.Background()
	userIDZero := helperCreateTest(t, randomTest(), 1)
	userID := helperCreateTest(t, randomTest(), 2)

	type args struct {
		repo   *RepositoryTests
		testID int
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
			name: "Success: delete test",
			args: args{
				repo:   mockRepo,
				testID: userID,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Fail: delete test with id = 0",
			args: args{
				repo:   mockRepo,
				testID: 0,
			},
			want: want{
				err: ErrTestAuthorIDEmpty,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.repo.DeleteTestById(ctx, tt.args.testID)
			if err != tt.want.err {
				t.Errorf("RepositoryTests.DeleteTestById() error = %v, wantErr %v", err, tt.want.err)
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteTest(t, userID)
		helperDeleteTest(t, userIDZero)
	})
}

func randomTest() domain.TestInput {
	return domain.TestInput{
		Title: util.RandomString(10),
	}
}

func helperDeleteTest(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE id = $1", id); err != nil {
		t.Errorf("error deleting test: %v", err)
	}
}

func helperCreateTest(t *testing.T, test domain.TestInput, authorID int) int {
	t.Helper()
	var id int
	if err := mockDB.QueryRow("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id", test.Title, authorID).Scan(&id); err != nil {
		t.Errorf("error creating test: %v", err)
	}
	return id
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
