package sessions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/popeskul/qna-go/internal/domain"
)

// RepositorySessions provides all the functions for the test repository.
type RepositorySessions struct {
	db *sql.DB
}

// NewRepoSessions creates a new instance of RepositoryTests.
func NewRepoSessions(db *sql.DB) *RepositorySessions {
	return &RepositorySessions{
		db: db,
	}
}

func (r *RepositorySessions) CreateRefreshToken(ctx context.Context, token domain.RefreshSession) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)", token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (r *RepositorySessions) GetRefreshToken(ctx context.Context, token string) (domain.RefreshSession, error) {
	queryGetRefreshToken := fmt.Sprintf("SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = %s", token)

	var t domain.RefreshSession
	err := r.db.QueryRowContext(ctx, queryGetRefreshToken).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = r.db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", t.UserID)

	return t, err
}
