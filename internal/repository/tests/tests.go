// Package tests is a struct that contains all functions for the test repository.
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
	ErrEmptyRepository   = errors.New("repository is empty")
)

// RepositoryTests provides all the functions to execute the queries and transactions.
type RepositoryTests struct {
	db *sql.DB
	*queries.Queries
}

// NewRepoTests creates a new instance of RepositoryTests.
func NewRepoTests(db *sql.DB) *RepositoryTests {
	return &RepositoryTests{
		db:      db,
		Queries: queries.NewQueries(db),
	}
}

// CreateTest creates a new test in the database.
func (r *RepositoryTests) CreateTest(ctx context.Context, authorID int, inputTest domain.TestInput) (int, error) {
	var id int

	err := r.ExecTx(ctx, func(tx *sql.Tx) error {
		if inputTest.Title == "" {
			return ErrEmptyTitle
		}

		if authorID == 0 {
			return ErrTestAuthorIDEmpty
		}

		createTestQuery := fmt.Sprintln("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id")
		if err := r.db.QueryRowContext(ctx, createTestQuery, inputTest.Title, authorID).Scan(&id); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

// GetTest returns a test by id.
// Returns the test and an error if any.
func (r *RepositoryTests) GetTest(ctx context.Context, testID int) (domain.Test, error) {
	var test domain.Test
	getTestQuery := fmt.Sprintln("SELECT * FROM tests WHERE id = $1")
	if err := r.db.QueryRowContext(ctx, getTestQuery, testID).Scan(&test.ID, &test.Title, &test.AuthorID, &test.CreatedAt, &test.UpdatedAt); err != nil {
		return domain.Test{}, err
	}

	return test, nil
}

// GetAllTestsByCurrentUser get all test from db by user id.
// It's return []domain.Test and error if any.
func (r *RepositoryTests) GetAllTestsByCurrentUser(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error) {
	allTests := make([]domain.Test, 0)
	allTestsQuery := fmt.Sprintln("SELECT * FROM tests WHERE author_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3")

	rows, err := r.db.QueryContext(ctx, allTestsQuery, userID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}

	var t domain.Test
	for rows.Next() {
		if err = rows.Scan(&t.ID, &t.Title, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		allTests = append(allTests, t)
	}
	err = rows.Err()

	return allTests, err
}

// UpdateTestById updates a test by id.
// Returns an error if any.
func (r *RepositoryTests) UpdateTestById(ctx context.Context, testID int, inputTest domain.TestInput) error {
	err := r.ExecTx(ctx, func(tx *sql.Tx) error {
		if inputTest.Title == "" {
			return ErrEmptyTitle
		}

		if testID == 0 {
			return ErrTestAuthorIDEmpty
		}

		updateTestQuery := fmt.Sprintln("UPDATE tests SET title = $1 WHERE id = $2")
		if _, err := r.db.ExecContext(ctx, updateTestQuery, inputTest.Title, testID); err != nil {
			return err
		}

		return nil
	})

	return err
}

// DeleteTestById deletes a test by id.
// Returns an error if any.
func (r *RepositoryTests) DeleteTestById(ctx context.Context, testID int) error {
	err := r.ExecTx(ctx, func(tx *sql.Tx) error {
		if testID == 0 {
			return ErrTestAuthorIDEmpty
		}

		deleteTestQuery := fmt.Sprintln("DELETE FROM tests WHERE id = $1")
		if _, err := r.db.ExecContext(ctx, deleteTestQuery, testID); err != nil {
			return err
		}

		return nil
	})

	return err
}
