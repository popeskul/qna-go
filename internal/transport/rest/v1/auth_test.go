package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func TestAuth_SignUp(t *testing.T) {
	ctx := context.Background()
	u := randomUser()
	u2 := randomUser()

	userID := helperCreatUser(t, ctx, u)

	validJSON, _ := json.Marshal(u2)
	invalidUniqueEmailJSON, _ := json.Marshal(u)
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewReader(tt.user))
			req.Header.Set("Content-Type", "application/json")

			r := gin.Default()
			r.POST("/api/v1/auth/sign-up", mockHandlers.SignUp)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				t.Cleanup(func() {
					var obj map[string]interface{}
					if err := json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
						t.Errorf("error unmarshalling response: %v", err)
					}

					// if user is created, and it has an id
					// then delete it
					if obj["id"] != nil {
						id := int(obj["id"].(float64))
						helperDeleteUserByID(t, id)
					}
				})

				return w.Code == tt.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
	})
}

func TestAuth_SignIn(t *testing.T) {
	ctx := context.Background()
	u := randomUser()

	userID := helperCreatUser(t, ctx, u)

	validJSON, _ := json.Marshal(u)
	invalidJSON, _ := json.Marshal(randomUser())
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(tt.user))
			req.Header.Set("Content-Type", "application/json")

			r := gin.Default()
			r.POST("/api/v1/auth/sign-in", mockHandlers.SignIn)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
	})
}
