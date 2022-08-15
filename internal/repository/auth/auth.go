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
	defer tx.Rollback()

	var id int
	createUserQuery := fmt.Sprintf("INSERT INTO users (name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id")
	if err = r.db.QueryRow(createUserQuery, u.Name, u.Email, u.Password).Scan(&id); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *RepositoryAuth) GetUser(email, password string) (domain.User, error) {
	var user domain.User

	getUserQuery := fmt.Sprintf("SELECT id, name, email, encrypted_password, created_at, updated_at FROM users WHERE email = $1 AND encrypted_password = $2")
	err := r.db.QueryRow(getUserQuery, email, password).Scan(&user.ID, &user.Name, &user.Email, &user.EncryptedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
