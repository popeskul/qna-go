// Package domain
// This place define test domain: Test.
package domain

// Test describe test entity.
type Test struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title" validate:"required,min=3,max=255"`
	AuthorID  int    `json:"author_id" db:"author_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

// ?
type GetAllTestsRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=10"`
}

// ?
type GetAllTestsParams struct {
	Limit  int `form:"limit" binding:"required,min=1"`
	Offset int `form:"offset" binding:"required,min=0"`
}
