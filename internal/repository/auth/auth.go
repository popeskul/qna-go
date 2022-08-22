// Package auth is a struct that contains all functions for the auth repository.
package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/queries"
)

// RepositoryAuth provides all the functions to execute the queries and transactions.
type RepositoryAuth struct {
	db *sql.DB
	*queries.Queries
}

// NewRepoAuth returns a new instance of the repository.
func NewRepoAuth(db *sql.DB) *RepositoryAuth {
	return &RepositoryAuth{
		db:      db,
		Queries: queries.NewQueries(db),
	}
}

// CreateUser creates a new user in the database.
// Returns the user and an error if any.
func (r *RepositoryAuth) CreateUser(ctx context.Context, u domain.SignUpInput) error {
	var userID int

	err := r.ExecTx(ctx, func(tx *sql.Tx) error {
		createUserQuery := fmt.Sprintln("INSERT INTO users (name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id")
		if err := r.db.QueryRowContext(ctx, createUserQuery, u.Name, u.Email, u.Password).Scan(&userID); err != nil {
			return err
		}

		return nil
	})

	return err
}

// GetUser returns a user from the database.
// Returns the user and an error if any.
func (r *RepositoryAuth) GetUser(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User

	getUserQuery := fmt.Sprintln("SELECT id, name, email, encrypted_password, created_at, updated_at FROM users WHERE email = $1 AND encrypted_password = $2")
	err := r.db.QueryRowContext(ctx, getUserQuery, email, password).Scan(&user.ID, &user.Name, &user.Email, &user.EncryptedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

// DeleteUserById deletes a user from the database.
// Returns an error if any.
func (r *RepositoryAuth) DeleteUserById(ctx context.Context, userID int) error {
	err := r.ExecTx(ctx, func(tx *sql.Tx) error {
		deleteUserQuery := fmt.Sprintln("DELETE FROM users WHERE id = $1")
		if _, err := r.db.ExecContext(ctx, deleteUserQuery, userID); err != nil {
			return err
		}

		return nil
	})

	return err
}
