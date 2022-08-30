// Package domain
// This place define basic auth domain: User.
package domain

// User describe user entity.
type User struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name" validate:"required,min=3,max=255"`
	Email     string `json:"email" db:"email" binding:"required" validate:"required,email,min=3,max=255"`
	Password  string `json:"password" db:"password" binding:"required" validate:"required,min=6,max=255"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
