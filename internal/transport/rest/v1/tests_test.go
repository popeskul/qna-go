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

	if err := mockServices.Auth.CreateUser(ctx, user); err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	userID, err := findUserIDByEmail(user.Email)
	if err != nil {
		t.Fatalf("error finding user: %v", err)
	}

	//token, err := mockServices.TokenMaker.CreateToken(userID, duration)
	token, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
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
						test := helperFindTestByTitle(t, test.Title)
						helperDeleteTestByID(t, test.ID)
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

func TestHandlers_GetTestByID(t *testing.T) {
	ctx := context.Background()
	user := randomUser()
	test := randomTest()

	if err := mockServices.Auth.CreateUser(ctx, user); err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	userID, err := findUserIDByEmail(user.Email)
	if err != nil {
		t.Fatalf("error finding user: %v", err)
	}
	test.AuthorID = userID

	//duration := time.Duration(1) * time.Second
	//token, err := mockServices.Manager.CreateToken(userID, duration)
	token, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	if err = mockRepo.CreateTest(ctx, user.ID, test); err != nil {
		t.Fatalf("error creating test: %v", err)
	}

	foundTest := helperFindTestByTitle(t, test.Title)

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
				id:    foundTest.ID,
				token: token,
			},
			want: want{
				status: http.StatusOK,
				test: domain.Test{
					ID:    foundTest.ID,
					Title: foundTest.Title,
				},
			},
		},
		{
			name: "Error: with invalid token",
			args: args{
				id:    foundTest.ID,
				token: "bad token",
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "No test found",
			args: args{
				id:    0,
				token: token,
			},
			want: want{
				status: http.StatusNotFound,
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
				trueStatus := w.Code == tt.want.status

				t.Log("->", w.Body.String())

				if w.Code == http.StatusOK {
					var test domain.Test
					if err = json.Unmarshal(w.Body.Bytes(), &test); err != nil {
						t.Fatalf("error unmarshalling test: %v", err)
					}

					return trueStatus && test.ID == tt.want.test.ID && test.Title == tt.want.test.Title
				}

				return trueStatus
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, user.ID)
		helperDeleteTestByID(t, foundTest.ID)
	})
}

func TestHandlers_GetAllTestsByCurrentUser(t *testing.T) {
	ctx := context.Background()
	user := randomUser()

	if err := mockServices.Auth.CreateUser(ctx, user); err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	userID, err := findUserIDByEmail(user.Email)
	if err != nil {
		t.Fatalf("error finding user: %v", err)
	}
	user.ID = userID

	token, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	type args struct {
		repo             *repository.Repository
		createByQuantity int
		params           domain.GetAllTestsRequest
		token            string
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
				token: token,
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
				token: token,
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
				token: token,
			},
			want: want{
				createdByQuantity: 2,
				status:            http.StatusOK,
			},
		},
		{
			name: "Fail: 0 tests with invalid pagination",
			args: args{
				repo:             mockRepo,
				createByQuantity: 10,
				params: domain.GetAllTestsRequest{
					PageID:   0,
					PageSize: -10,
				},
				token: token,
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
			req.Header.Set("Authorization", "Bearer "+tt.args.token)

			q := req.URL.Query()
			q.Add("page_id", strconv.Itoa(tt.args.params.PageID))
			q.Add("page_size", strconv.Itoa(tt.args.params.PageSize))
			req.URL.RawQuery = q.Encode()

			r := gin.Default()
			r.GET("/api/v1/tests", mockHandlers.authMiddleware, mockHandlers.GetAllTestsByUserID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				t.Log(w.Body.String())
				if w.Code != tt.want.status {
					t.Fatalf("error getting tests: %v", w.Code)
				}
				if w.Code == http.StatusBadRequest {
					return true
				}

				var tests []domain.Test
				if err := json.Unmarshal(w.Body.Bytes(), &tests); err != nil {
					t.Errorf("error unmarshalling tests: %v", err)
				}

				if len(tests) != tt.want.createdByQuantity {
					t.Errorf("got %v, want %v", len(tests), tt.want.createdByQuantity)
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

	if err := mockServices.Auth.CreateUser(ctx, user); err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	userID, err := findUserIDByEmail(user.Email)
	if err != nil {
		t.Fatalf("error finding user id: %v", err)
	}

	//duration := time.Duration(1) * time.Second
	//token, err := mockServices.Manager.CreateToken(userID, duration)
	token, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}
	if err != nil {
		t.Fatalf("error generating token: %v", err)
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

	if err := mockServices.Auth.CreateUser(ctx, user); err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	userID, err := findUserIDByEmail(user.Email)
	if err != nil {
		t.Fatalf("error finding user id: %v", err)
	}

	//duration := time.Duration(1) * time.Second
	//token, err := mockServices.Manager.CreateToken(userID, duration)
	token, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	testID1 := helperCreateTest(t, userID, randomTest())
	testID2 := helperCreateTest(t, userID, randomTest())

	type args struct {
		token string
		id    int
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
				id:    testID2,
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "Fail: Delete test that does not exist",
			args: args{
				token: token,
				id:    2134234234,
			},
			want: want{
				status: http.StatusNotFound,
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
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/tests/"+strconv.Itoa(tt.args.id), nil)
			req.Header.Set("Authorization", "Bearer "+tt.args.token)

			r := gin.Default()
			r.DELETE("/api/v1/tests/:id", mockHandlers.authMiddleware, mockHandlers.DeleteTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.want.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
		helperDeleteTestByID(t, testID1)
		helperDeleteTestByID(t, testID2)
	})
}
