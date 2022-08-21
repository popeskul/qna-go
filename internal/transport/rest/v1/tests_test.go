package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/util"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandlers_CreateTests(t *testing.T) {
	ctx := context.Background()
	user := randomUser()
	test := randomTest()

	userID, err := mockServices.Auth.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(ctx, user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

	validJSON, _ := json.Marshal(test)
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
						if err = json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
							t.Errorf("error unmarshalling response: %v", err)
						}

						if obj["id"] != nil {
							testID := int(obj["id"].(float64))
							helperDeleteTestByID(t, testID)
						}
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

func TestHandlers_GetTest(t *testing.T) {
	ctx := context.Background()
	user := randomUser()
	test := randomTest()

	userID, err := mockServices.Auth.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(ctx, user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

	testID, err := mockRepo.CreateTest(ctx, userID, test)
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
				if err = json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
					t.Errorf("error unmarshalling response: %v", err)
				}

				resTest := obj["test"].(map[string]interface{})
				if id := int(resTest["id"].(float64)); id != tt.want.test.ID {
					t.Errorf("error getting test: %v", resTest["id"])
				}

				resTestID := int(resTest["id"].(float64))
				resTestTitle := resTest["title"].(string)

				if resTestTitle != tt.want.test.Title {
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

func TestHandlers_GetAllTestsByUserID(t *testing.T) {
	ctx := context.Background()
	user := randomUser()

	userID, err := mockServices.Auth.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(ctx, user.Email, user.Password)

	type args struct {
		repo             *repository.Repository
		createByQuantity int
		params           domain.GetAllTestsRequest
	}
	type want struct {
		createdByQuantity int
		status            int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Success: Get all tests with default pagination",
			args: args{
				repo:             mockRepo,
				createByQuantity: 10,
				params: domain.GetAllTestsRequest{
					PageID:   1,
					PageSize: 10,
				},
			},
			want: want{
				createdByQuantity: 10,
				status:            http.StatusOK,
			},
		},
		{
			name: "Success: Get 1 page of tests but in db there are more than 10 tests",
			args: args{
				repo:             mockRepo,
				createByQuantity: 12,
				params: domain.GetAllTestsRequest{
					PageID:   1,
					PageSize: 10,
				},
			},
			want: want{
				createdByQuantity: 10,
				status:            http.StatusOK,
			},
		},
		{
			name: "Success: Get 3 page with 2 tests, in db there are 22 tests",
			args: args{
				repo:             mockRepo,
				createByQuantity: 22,
				params: domain.GetAllTestsRequest{
					PageID:   3,
					PageSize: 10,
				},
			},
			want: want{
				createdByQuantity: 2,
				status:            http.StatusOK,
			},
		},
		{
			name: "Fail: Get all tests with invalid pagination",
			args: args{
				repo:             mockRepo,
				createByQuantity: 10,
				params: domain.GetAllTestsRequest{
					PageID:   0,
					PageSize: -10,
				},
			},
			want: want{
				createdByQuantity: 0,
				status:            http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdIDs := make([]int, tt.args.createByQuantity)

			for i := 0; i < tt.args.createByQuantity; i++ {
				testID := helperCreateTest(t, userID, randomTest())
				createdIDs = append(createdIDs, testID)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/tests", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			q := req.URL.Query()
			q.Add("page_id", strconv.Itoa(tt.args.params.PageID))
			q.Add("page_size", strconv.Itoa(tt.args.params.PageSize))
			req.URL.RawQuery = q.Encode()

			r := gin.Default()
			r.GET("/api/v1/tests", mockHandlers.authMiddleware, mockHandlers.GetAllTestsByCurrentUser)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				if w.Code != tt.want.status {
					t.Errorf("error getting tests: %v", w.Code)
				} else {
					// we need to exit from test if status is not ok, and it's by plan
					if w.Code == http.StatusBadRequest {
						return true
					}

					res := w.Body.String()
					obj := make(map[string]interface{})
					t.Log("response: ", res)
					err = json.Unmarshal([]byte(res), &obj)
					if err != nil {
						t.Errorf("error unmarshalling response: %v", err)
					}

					tests := obj["tests"].([]interface{})
					if len(tests) != tt.want.createdByQuantity {
						t.Errorf("error getting tests: %v", len(tests))
					}
				}

				return w.Code == tt.want.status
			})

			t.Cleanup(func() {
				for _, testID := range createdIDs {
					helperDeleteTestByID(t, testID)
				}
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
	})
}

func TestHandlers_UpdateTestByID(t *testing.T) {
	ctx := context.Background()
	user := randomUser()

	userID, err := mockServices.Auth.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(ctx, user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

	newTitle := util.RandomString(10)
	validJSON := []byte(`{"title": "` + newTitle + `"}`)
	badJSON := []byte(`bad request`)

	testIDZero := helperCreateTest(t, userID, randomTest())
	testID := helperCreateTest(t, userID, randomTest())

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
	ctx := context.Background()
	user := randomUser()

	userID, err := mockServices.Auth.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	token, err := mockServices.Auth.GenerateToken(ctx, user.Email, user.Password)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}

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
