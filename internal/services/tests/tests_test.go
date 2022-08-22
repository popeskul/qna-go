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
	"testing"

	_ "github.com/lib/pq"
)

var (
	mockDB   *sql.DB
	mockRepo *repository.Repository
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
			err := mockRepo.CreateTest(ctx, tt.args.userID, tt.args.input)

			if err != tt.want.err {
				t.Fatalf("ServiceTests.CreateTest() error = %v, wantErr %v", err, tt.want.err)
			}

			t.Cleanup(func() {
				helperDeleteTestByTitle(t, tt.args.input.Title)
			})
		})
	}
}

func TestServiceTests_UpdateTestById(t *testing.T) {
	ctx := context.Background()
	mockUserID := 1
	test := randomTest()
	if err := mockRepo.CreateTest(ctx, mockUserID, test); err != nil {
		t.Fatalf("Some error occured. Err: %s", err)
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
			err := mockRepo.UpdateTestById(ctx, tt.args.userID, tt.args.input)

			if err != tt.want.err {
				t.Fatalf("ServiceTests.UpdateTestById() error = %v, wantErr %v", err, tt.want.err)
			}
		})
	}

	t.Cleanup(func() {
		helperDeleteTestByTitle(t, test.Title)
	})
}

func TestServiceTests_DeleteTestById(t *testing.T) {
	ctx := context.Background()
	mockUserID := 1

	testMock := randomTest()
	if err := mockRepo.CreateTest(ctx, mockUserID, testMock); err != nil {
		t.Fatalf("Some error occured. Err: %s", err)
	}

	test := helperFindTestByTitle(t, testMock.Title)

	type args struct {
		repo   *repository.Repository
		testID int
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
				testID: test.ID,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Fail: Delete test",
			args: args{
				repo:   mockRepo,
				testID: 11111111,
			},
			want: want{
				err: tests.ErrDeleteTest,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := mockRepo.DeleteTestById(ctx, tt.args.testID)
			if err != tt.want.err {
				t.Fatalf("ServiceTests.DeleteTestById() error = %v, wantErr %v", err, tt.want.err)
			}
		})
	}
}

func TestServiceTests_GetAllTests(t *testing.T) {
	ctx := context.Background()
	userID := 1

	type args struct {
		repo             *repository.Repository
		CreateByQuantity int
		params           domain.GetAllTestsParams
	}
	type want struct {
		CreatedByQuantity int
	}

	tests := []struct {
		name string
		args struct {
			repo             *repository.Repository
			CreateByQuantity int
			params           domain.GetAllTestsParams
		}
		want struct {
			CreatedByQuantity int
		}
	}{
		{
			name: "Success: get all tests (3)",
			args: args{
				repo:             mockRepo,
				CreateByQuantity: 3,
				params: domain.GetAllTestsParams{
					Limit:  10,
					Offset: 0,
				},
			},
			want: want{
				CreatedByQuantity: 3,
			},
		},
		{
			name: "Success: get all tests from empty DB",
			args: args{
				repo:             mockRepo,
				CreateByQuantity: 0,
			},
			want: want{
				CreatedByQuantity: 0,
			},
		},
		{
			name: "Success: get all tests from DB with limit",
			args: args{
				repo:             mockRepo,
				CreateByQuantity: 5,
				params: domain.GetAllTestsParams{
					Limit:  3,
					Offset: 4,
				},
			},
			want: want{
				CreatedByQuantity: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdIDs := make([]int, tt.args.CreateByQuantity)

			for i := 0; i < tt.args.CreateByQuantity; i++ {
				testID := helperCreateTest(t, userID, randomTest())
				createdIDs = append(createdIDs, testID)
			}

			allTests, err := tt.args.repo.GetAllTestsByCurrentUser(ctx, userID, tt.args.params)
			if err != nil {
				t.Errorf("RepositoryTests.GetAllTests() error = %v", err)
			}

			// check count of tests
			if len(allTests) != tt.want.CreatedByQuantity {
				t.Errorf("RepositoryTests.GetAllTests() error = %v, wantErr %v", len(allTests), tt.want.CreatedByQuantity)
			}

			t.Cleanup(func() {
				for _, testID := range createdIDs {
					t.Log("Delete test: ", testID)
					helperDeleteTest(t, testID)
				}
			})
		})
	}
}

func randomTest() domain.TestInput {
	return domain.TestInput{
		Title: util.RandomString(10),
	}
}

func helperCreateTest(t *testing.T, authorID int, test domain.TestInput) int {
	t.Helper()
	var id int
	if err := mockDB.QueryRow("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id", test.Title, authorID).Scan(&id); err != nil {
		t.Fatalf("helperCreateTest() error = %v", err)
	}
	return id
}

func helperDeleteTest(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE id = $1", id); err != nil {
		t.Errorf("error deleting test: %v", err)
	}
}

func helperFindTestByTitle(t *testing.T, title string) domain.Test {
	t.Helper()
	var test domain.Test
	if err := mockDB.QueryRow("SELECT id, title, author_id FROM tests WHERE title = $1", title).Scan(&test.ID, &test.Title, &test.AuthorID); err != nil {
		t.Fatalf("helperFindTestByTitle() error = %v", err)
	}
	return test
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
