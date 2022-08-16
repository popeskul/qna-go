package auth

import (
	"database/sql"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
)

type RepositoryAuth struct {
	db *sql.DB
}

func NewRepoAuth(db *sql.DB) *RepositoryAuth {
	return &RepositoryAuth{
		db: db,
	}
}

func (r *RepositoryAuth) CreateUser(u domain.SignUpInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println("failed to rollback transaction: ", err)
		}
	}(tx)

	var id int
	createUserQuery := fmt.Sprintln("INSERT INTO users (name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id")
	if err = r.db.QueryRow(createUserQuery, u.Name, u.Email, u.Password).Scan(&id); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *RepositoryAuth) GetUser(email, password string) (domain.User, error) {
	var user domain.User

	getUserQuery := fmt.Sprintln("SELECT id, name, email, encrypted_password, created_at, updated_at FROM users WHERE email = $1 AND encrypted_password = $2")
	err := r.db.QueryRow(getUserQuery, email, password).Scan(&user.ID, &user.Name, &user.Email, &user.EncryptedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *RepositoryAuth) DeleteUserById(userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println("failed to rollback transaction: ", err)
		}
	}(tx)

	deleteUserQuery := fmt.Sprintln("DELETE FROM users WHERE id = $1")
	if _, err := r.db.Exec(deleteUserQuery, userID); err != nil {
		return err
	}

	return tx.Commit()
}
