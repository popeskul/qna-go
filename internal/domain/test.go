// Package domain
// This place define test domain: Test, TestInput.
package domain

// Test represents response for the test
type Test struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	AuthorID  int    `json:"author_id" db:"author_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

// TestInput represents request for the test
type TestInput struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
}
