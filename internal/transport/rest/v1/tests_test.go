package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandlers_CreateTests(t *testing.T) {
	user := domain.SignUpInput{
		Email:    "TestHandlers_CreateTests@mail.com",
		Password: "TestHandlers_CreateTests",
	}

	test := domain.TestInput{
		Title: "TestHandlers_CreateTests",
	}

	userID, err := mockServices.Auth.CreateUser(user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

	validJSON := []byte(`{"title": "` + test.Title + `"}`)
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
				t.Cleanup(func() {
					if w.Code == http.StatusCreated {
						var obj map[string]interface{}
						if err := json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
							t.Errorf("error unmarshalling response: %v", err)
						}

						testID := int(obj["id"].(float64))

						helperDeleteTestByID(t, testID)
					}
				})
				return w.Code == tt.status

			})
		})
	}

	t.Cleanup(func() {
		fmt.Println("cleanup")
		helperDeleteUserByID(t, userID)
	})
}

func TestHandlers_GetTest(t *testing.T) {
	user := domain.SignUpInput{
		Email:    "TestHandlers_GetTest@mail.com",
		Password: "1231",
	}
	userID, err := mockServices.Auth.CreateUser(user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

	test := domain.TestInput{
		Title: "TestHandlers_GetTest",
	}

	testID, err := mockRepo.CreateTest(1, test)
	if err != nil {
		t.Errorf("error creating test: %v", err)
	}

	type args struct {
		id    int
		token string
	}
	type want struct {
		status int
		test   domain.Test
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: Get test",
			args: args{
				id:    testID,
				token: token,
			},
			want: want{
				status: http.StatusOK,
				test: domain.Test{
					ID:    testID,
					Title: test.Title,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tests/"+strconv.Itoa(tt.args.id), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.args.token)

			r := gin.Default()
			r.GET("/api/v1/tests/:id", mockHandlers.authMiddleware, mockHandlers.GetTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				var obj map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
					t.Errorf("error unmarshalling response: %v", err)
				}

				status := obj["status"].(string)
				if status != "success" {
					t.Errorf("error getting test: %v", status)
				}

				resTest := obj["test"].(map[string]interface{})
				resTestID := int(resTest["id"].(float64))
				resTestTitle := resTest["title"].(string)
				fmt.Println("resTestID: ", resTestID, "resTestTitle: ", resTestTitle)

				if testID != tt.want.test.ID {
					t.Errorf("error getting test: %v", testID)
				}

				if test.Title != resTestTitle {
					t.Errorf("error getting test: %v", test.Title)
				}

				t.Cleanup(func() {
					if w.Code == http.StatusOK {
						helperDeleteTestByID(t, resTestID)
					}
				})

				return w.Code == tt.want.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
		helperDeleteTestByID(t, testID)
	})
}

func TestHandlers_UpdateTestByID(t *testing.T) {
	user := domain.SignUpInput{
		Email:    "ad2@mail.com",
		Password: "ad",
	}

	userID, err := mockServices.Auth.CreateUser(user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(user.Email, user.Password)

	newTitle := "title1"
	validJSON := []byte(`{"title": "` + newTitle + `"}`)
	badJSON := []byte(`bad request`)

	testIDZero := helperCreateTest(t, userID, domain.TestInput{Title: "title0"})
	testID := helperCreateTest(t, userID, domain.TestInput{Title: "title1"})

	type args struct {
		id    int
		input []byte
		token string
	}
	type want struct {
		title  string
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: Update test",
			args: args{
				token: token,
				id:    testID,
				input: validJSON,
			},
			want: want{
				title:  newTitle,
				status: http.StatusOK,
			},
		},
		{
			name: "Error: with invalid json",
			args: args{
				token: token,
				input: badJSON,
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Error: invalid token",
			args: args{
				id:    testID,
				input: validJSON,
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/v1/tests/"+strconv.Itoa(tt.args.id), bytes.NewReader(tt.args.input))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.args.token)

			r := gin.Default()
			r.PUT("/api/v1/tests/:id", mockHandlers.authMiddleware, mockHandlers.UpdateTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.want.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteTestByID(t, testIDZero)
		helperDeleteTestByID(t, testID)
		helperDeleteUserByID(t, userID)
	})
}

func TestHandlers_DeleteTestByID(t *testing.T) {
	user := domain.SignUpInput{
		Email:    "TestHandlers_DeleteTestByID@mail.com",
		Password: "12345",
	}

	userID, err := mockServices.Auth.CreateUser(user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(user.Email, user.Password)

	type args struct {
		token string
	}
	type want struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: Delete test",
			args: args{
				token: token,
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "Error: invalid token",
			args: args{},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testID := helperCreateTest(t, userID, domain.TestInput{Title: "title"})
			if testID == 0 {
				t.Errorf("error creating test: %v", testID)
			}

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/tests/"+strconv.Itoa(testID), nil)
			req.Header.Set("Authorization", "Bearer "+tt.args.token)

			r := gin.Default()
			r.DELETE("/api/v1/tests/:id", mockHandlers.authMiddleware, mockHandlers.DeleteTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.want.status
			})

			t.Cleanup(func() {
				helperDeleteTestByID(t, testID)
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
	})
}

func helperCreateTest(t *testing.T, userID int, test domain.TestInput) int {
	t.Helper()

	var id int
	if err := mockDB.QueryRow("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id", test.Title, userID).Scan(&id); err != nil {
		t.Errorf("error inserting test: %v", err)
	}

	return id
}

func helperDeleteTestByID(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE id = $1", id); err != nil {
		t.Errorf("error deleting test: %v", err)
	}
}

func helperDeleteUserByID(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		t.Errorf("error deleting user: %v", err)
	}
}
