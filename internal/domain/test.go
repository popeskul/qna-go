package domain

type Test struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	AuthorID  int    `json:"author_id" db:"author_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type TestInput struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
}
