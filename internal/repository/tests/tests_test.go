package tests

import (
	"github.com/popeskul/qna-go/internal/domain"
	"testing"
)

func TestRepositoryTests_CreateTest(t *testing.T) {
	mockUniqueTest := domain.TestInput{
		Title: "CreateTest 1",
	}
	mockUserID := 1

	uniqueID := helperCreateTest(t, mockUniqueTest.Title, mockUserID)

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
				repo: mockRepo,
				input: domain.TestInput{
					Title: mockUniqueTest.Title,
				},
				userID: mockUserID,
			},
			want: want{
				id:  uniqueID + 1,
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
			id, err := tt.args.repo.CreateTest(tt.args.userID, tt.args.input)

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
		helperDeleteTest(t, uniqueID)
	})
}

func TestRepositoryTests_UpdateTestById(t *testing.T) {
	mockTestTitle := "Repo UpdateTestById 1"
	mockTestAuthorID := 1

	createdID := helperCreateTest(t, mockTestTitle, mockTestAuthorID)

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
			err := tt.args.repo.UpdateTestById(tt.args.testID, tt.args.input)
			if err != tt.want.err {
				t.Errorf("RepositoryTests.UpdateTestById() error = %v, wantErr %v", err, tt.want.err)
			}

			test, _ := tt.args.repo.GetTest(createdID)
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
	userIDZero := helperCreateTest(t, "Repo DeleteTestById 1", 1)
	userID := helperCreateTest(t, "DeleteTestById 2", 2)

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
			err := tt.args.repo.DeleteTestById(tt.args.testID)
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

func helperDeleteTest(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE id = $1", id); err != nil {
		t.Errorf("error deleting test: %v", err)
	}
}

func helperCreateTest(t *testing.T, title string, authorID int) int {
	t.Helper()
	var id int
	if err := mockDB.QueryRow("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id", title, authorID).Scan(&id); err != nil {
		t.Errorf("error creating test: %v", err)
	}
	return id
}
