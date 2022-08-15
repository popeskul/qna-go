package v1

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func TestAuth_SignUp(t *testing.T) {
	mockUniqEmail := "TestAuth_SignUpUnique@mail.com"
	mockEmail := "TestAuth_SignUp@mail.com"
	mockPassword := "123456"

	validJSON := []byte(`{"email": "` + mockEmail + `", "encrypted_password": "` + mockPassword + `", "name": "1123"}`)
	invalidUniqueEmailJSON := []byte(`{"email": "` + mockUniqEmail + `", "encrypted_password": "` + mockPassword + `", "name": "2123"}`)
	badJSON := []byte(`bad request`)

	tests := []struct {
		name   string
		user   []byte
		status int
	}{
		{
			name:   "Success: Create user",
			user:   validJSON,
			status: http.StatusCreated,
		},
		{
			name:   "Error: with existing email",
			user:   invalidUniqueEmailJSON,
			status: http.StatusInternalServerError,
		},
		{
			name:   "Error: with invalid json",
			user:   badJSON,
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockDB.Exec("INSERT INTO users (email, encrypted_password, name) VALUES ($1, $2, $3) RETURNING id", mockUniqEmail, mockPassword, "-"); err != nil {
				t.Errorf("Some error occured. Err: %mockServices", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewReader(tt.user))
			req.Header.Set("Content-Type", "application/json")

			r := gin.Default()
			r.POST("/api/v1/auth/sign-up", mockHandlers.SignUp)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.status
			})

			t.Cleanup(func() {
				if _, err := mockDB.Exec("DELETE FROM users WHERE email IN ($1, $2)", mockEmail, mockUniqEmail); err != nil {
					t.Errorf("Some error occured. Err: %mockServices", err)
				}
			})
		})
	}
}

func TestAuth_SignIn(t *testing.T) {
	mockEmail := "TestAuth_SignUp@mail.com"
	mockPassword := "123456"

	validJSON := []byte(`{"email": "` + mockEmail + `", "password": "` + mockPassword + `"}`)
	invalidJSON := []byte(`{"email": "wrong@mail.com", "password": "` + mockPassword + `"}`)
	badJSON := []byte(`bad request`)

	tests := []struct {
		name   string
		user   []byte
		status int
	}{
		{
			name:   "Success: Sign in",
			user:   validJSON,
			status: http.StatusOK,
		},
		{
			name:   "Error: with invalid json",
			user:   badJSON,
			status: http.StatusBadRequest,
		},
		{
			name:   "Error: with invalid email",
			user:   invalidJSON,
			status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockServices.CreateUser(domain.SignUpInput{
				Email:    mockEmail,
				Password: mockPassword,
			}); err != nil {
				t.Errorf("Some error occured. Err: %mockServices", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(tt.user))
			req.Header.Set("Content-Type", "application/json")

			r := gin.Default()
			r.POST("/api/v1/auth/sign-in", mockHandlers.SignIn)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.status
			})

			t.Cleanup(func() {
				if _, err := mockDB.Exec("DELETE FROM users WHERE email IN ($1)", mockEmail); err != nil {
					t.Errorf("Some error occured. Err: %mockServices", err)
				}
			})
		})
	}
}
