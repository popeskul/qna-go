package tests

import (
	"github.com/popeskul/qna-go/internal/domain"
	"testing"
)

func TestRepositoryTests_CreateTest(t *testing.T) {
	mockUniqueTest := domain.TestInput{
		Title: "Test title",
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
