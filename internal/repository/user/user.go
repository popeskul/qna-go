// Package user is a struct that contains all functions for the auth repository.
package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
)

var (
	ErrCreateUser = errors.New("error creating user")
	ErrDeleteUser = errors.New("error deleting user")
)

// RepositoryAuth provides all the functions to execute the queries and transactions.
type RepositoryAuth struct {
	db *sql.DB
}

// NewRepoAuth returns a new instance of the repository.
func NewRepoAuth(db *sql.DB) *RepositoryAuth {
	return &RepositoryAuth{
		db: db,
	}
}

// CreateUser creates a new user in the database and an error if any.
func (r *RepositoryAuth) CreateUser(ctx context.Context, u domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	createUserQuery := fmt.Sprintln("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)")
	result, err := r.db.ExecContext(ctx, createUserQuery, u.Name, u.Email, u.Password)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrCreateUser
	}

	return tx.Commit()
}

// GetUser returns a user from the database and an error if any.
func (r *RepositoryAuth) GetUser(ctx context.Context, email string, password []byte) (domain.User, error) {
	var user domain.User

	getUserQuery := fmt.Sprintln("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND password = $2")
	err := r.db.QueryRowContext(ctx, getUserQuery, email, password).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByEmail returns a user from the database and an error if any.
func (r *RepositoryAuth) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	getUserQuery := fmt.Sprintln("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")
	err := r.db.QueryRowContext(ctx, getUserQuery, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

// DeleteUserById deletes a user from the database and returns an error if any.
func (r *RepositoryAuth) DeleteUserById(ctx context.Context, userID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	deleteUserQuery := fmt.Sprintln("DELETE FROM users WHERE id = $1")
	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrDeleteUser
	}

	return tx.Commit()
}
