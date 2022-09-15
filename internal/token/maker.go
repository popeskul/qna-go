package token

import "time"

// Manager is an interface for managing token.
type Manager interface {
	// CreateToken creates a new token for specific username and duration.
	CreateToken(userID int, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
