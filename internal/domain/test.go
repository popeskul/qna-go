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

type GetTestByIDRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type GetTestByIDResponse struct {
	Test Test `json:"test"`
}

type GetAllTestsRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=10"`
}

type AllTestResponse struct {
	Tests []Test `json:"tests"`
}

type GetAllTestsParams struct {
	Limit  int `form:"limit" binding:"required,min=1"`
	Offset int `form:"offset" binding:"required,min=0"`
}
