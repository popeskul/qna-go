package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/queries"
)

var (
	ErrDuplicateAuthorID = errors.New("pq: duplicate key value violates unique constraint \"tests_author_id_key\"")
	ErrTestAuthorIDEmpty = errors.New("author id is empty")
	ErrEmptyTitle        = errors.New("title is empty")
)

type RepositoryTests struct {
	db *sql.DB
	*queries.Queries
}

func NewRepoTests(db *sql.DB) *RepositoryTests {
	return &RepositoryTests{
		db:      db,
		Queries: queries.NewQueries(db),
	}
}

func (r *RepositoryTests) CreateTest(authorID int, inputTest domain.TestInput) (int, error) {
	var id int

	err := r.ExecTx(context.Background(), func(tx *sql.Tx) error {
		if inputTest.Title == "" {
			return ErrEmptyTitle
		}

		if authorID == 0 {
			return ErrTestAuthorIDEmpty
		}

		createTestQuery := fmt.Sprintln("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id")
		if err := r.db.QueryRow(createTestQuery, inputTest.Title, authorID).Scan(&id); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func (r *RepositoryTests) GetTest(testID int) (domain.Test, error) {
	var test domain.Test
	getTestQuery := fmt.Sprintln("SELECT * FROM tests WHERE id = $1")
	if err := r.db.QueryRow(getTestQuery, testID).Scan(&test.ID, &test.Title, &test.AuthorID, &test.CreatedAt, &test.UpdatedAt); err != nil {
		return domain.Test{}, err
	}

	return test, nil
}

func (r *RepositoryTests) UpdateTestById(testID int, inputTest domain.TestInput) error {
	err := r.ExecTx(context.Background(), func(tx *sql.Tx) error {
		if inputTest.Title == "" {
			return ErrEmptyTitle
		}

		if testID == 0 {
			return ErrTestAuthorIDEmpty
		}

		updateTestQuery := fmt.Sprintln("UPDATE tests SET title = $1 WHERE id = $2")
		if _, err := r.db.Exec(updateTestQuery, inputTest.Title, testID); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *RepositoryTests) DeleteTestById(testID int) error {
	err := r.ExecTx(context.Background(), func(tx *sql.Tx) error {
		if testID == 0 {
			return ErrTestAuthorIDEmpty
		}

		deleteTestQuery := fmt.Sprintln("DELETE FROM tests WHERE id = $1")
		if _, err := r.db.Exec(deleteTestQuery, testID); err != nil {
			return err
		}

		return nil
	})

	return err
}
