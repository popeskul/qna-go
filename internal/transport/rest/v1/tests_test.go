package v1

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_CreateTests(t *testing.T) {
	user := domain.SignUpInput{
		Email:    "TestHandlers_CreateTests@mail.com",
		Password: "TestHandlers_CreateTests",
	}

	_, err := mockServices.Auth.CreateUser(user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

	testTitle := "title1"

	validJSON := []byte(`{"title": "` + testTitle + `"}`)
	badJSON := []byte(`bad request`)

	tests := []struct {
		name   string
		test   []byte
		status int
	}{
		{
			name:   "Success: Create test",
			test:   validJSON,
			status: http.StatusCreated,
		},
		{
			name:   "Error: with invalid json",
			test:   badJSON,
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tests", bytes.NewReader(tt.test))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			r := gin.Default()
			r.POST("/api/v1/tests", mockHandlers.authMiddleware, mockHandlers.CreateTest)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.status
			})

			t.Cleanup(func() {
				helperDeleteTestByTitle(t, testTitle)
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByEmail(t, user.Email)
	})
}

func helperDeleteTestByTitle(t *testing.T, title string) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE title = $1", title); err != nil {
		t.Errorf("error deleting test: %v", err)
	}
}

func helperDeleteUserByEmail(t *testing.T, email string) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM users WHERE email = $1", email); err != nil {
		t.Errorf("error deleting user: %v", err)
	}
}
