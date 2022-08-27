// Package tests is a struct that contains all functions for the test repository.
package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
)

var (
	ErrDeleteTest   = errors.New("error deleting test")
	ErrTestNotFound = errors.New("test not found")
)

// RepositoryTests provides all the functions for the test repository.
type RepositoryTests struct {
	db *sql.DB
}

// NewRepoTests creates a new instance of RepositoryTests.
func NewRepoTests(db *sql.DB) *RepositoryTests {
	return &RepositoryTests{
		db: db,
	}
}

// CreateTest creates a new test in the database.
func (r *RepositoryTests) CreateTest(ctx context.Context, authorID int, inputTest domain.Test) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	var id int
	createTestQuery := fmt.Sprintln("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id")
	if err = r.db.QueryRowContext(ctx, createTestQuery, inputTest.Title, authorID).Scan(&id); err != nil {
		return err
	}

	return tx.Commit()
}

// GetTest returns a test by id and returns test and error if any.
func (r *RepositoryTests) GetTest(ctx context.Context, testID int) (domain.Test, error) {
	var test domain.Test
	getTestQuery := fmt.Sprintln("SELECT * FROM tests WHERE id = $1")
	if err := r.db.QueryRowContext(ctx, getTestQuery, testID).Scan(&test.ID, &test.Title, &test.AuthorID, &test.CreatedAt, &test.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return test, ErrTestNotFound
		}

		return domain.Test{}, err
	}

	return test, nil
}

// GetAllTestsByUserID get all test from db by user id and returns tests and error if any.
func (r *RepositoryTests) GetAllTestsByUserID(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error) {
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

// UpdateTestById updates a test by id and returns error if any.
func (r *RepositoryTests) UpdateTestById(ctx context.Context, testID int, inputTest domain.Test) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	updateTestQuery := fmt.Sprintln("UPDATE tests SET title = $1 WHERE id = $2")
	if _, err = r.db.ExecContext(ctx, updateTestQuery, inputTest.Title, testID); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteTestById deletes a test by id and returns error if any.
func (r *RepositoryTests) DeleteTestById(ctx context.Context, testID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	deleteTestQuery := fmt.Sprintln("DELETE FROM tests WHERE id = $1 RETURNING id")
	res, err := r.db.ExecContext(ctx, deleteTestQuery, testID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteTest
	}

	return tx.Commit()
}
