// Package domain
// This place define basic auth domain: User, SignUpInput, SignInInput.
package domain

// User represents response for the user
type User struct {
	ID                int    `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Email             string `json:"email" db:"email"`
	EncryptedPassword string `json:"encrypted_password" db:"encrypted_password"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	UpdatedAt         string `json:"updated_at" db:"updated_at"`
}

// SignUpInput represents request for the user sign up
type SignUpInput struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// SignInInput represents request for the user sign in
type SignInInput struct {
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}
