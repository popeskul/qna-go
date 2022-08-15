package tests

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
)

var (
	ErrDuplicateAuthorID = errors.New("pq: duplicate key value violates unique constraint \"tests_author_id_key\"")
	ErrTestTitleEmpty    = errors.New("test title is empty")
	ErrTestAuthorIDEmpty = errors.New("author id is empty")
)

type RepositoryTests struct {
	db *sql.DB
}

func NewRepoTests(db *sql.DB) *RepositoryTests {
	return &RepositoryTests{
		db: db,
	}
}

func (r *RepositoryTests) CreateTest(authorID int, inputTest domain.TestInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if inputTest.Title == "" {
		return 0, ErrTestTitleEmpty
	}

	if authorID == 0 {
		return 0, ErrTestAuthorIDEmpty
	}

	var id int
	createTestQuery := fmt.Sprintf("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id")
	if err = r.db.QueryRow(createTestQuery, inputTest.Title, authorID).Scan(&id); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *RepositoryTests) GetTest(testID int) (domain.Test, error) {
	var test domain.Test
	getTestQuery := fmt.Sprintf("SELECT * FROM tests WHERE id = $1")
	if err := r.db.QueryRow(getTestQuery, testID).Scan(&test.ID, &test.Title, &test.AuthorID, &test.CreatedAt, &test.UpdatedAt); err != nil {
		return domain.Test{}, err
	}

	return test, nil
}
