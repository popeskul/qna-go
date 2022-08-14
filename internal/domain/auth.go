package domain

type User struct {
	ID                int    `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Email             string `json:"email" db:"email"`
	EncryptedPassword string `json:"encrypted_password" db:"encrypted_password"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	UpdatedAt         string `json:"updated_at" db:"updated_at"`
}

type SignUpInput struct {
	Name              string `json:"name" validate:"required,min=3,max=255"`
	Email             string `json:"email" validate:"required,email,min=3,max=255"`
	EncryptedPassword string `json:"encrypted_password" validate:"required,min=6,max=255"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}
