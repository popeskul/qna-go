package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-contrib/sessions"
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

	accessToken, refreshToken, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
	}

	validJSON, _ := json.Marshal(test)
	badJSON := []byte(`bad request`)

	tests := []struct {
		name   string
		test   []byte
		status int
	}{
		{
			name:   "Success: CreateRefreshToken test",
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

			r := gin.Default()
			r.Use(sessions.Sessions("session", mockHandlers.store))

			r.POST("/api/v1/tests", setSessionMiddleware(t, accessToken), mockHandlers.authMiddleware, mockHandlers.CreateTest)

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
		helperDeleteRefreshTokenByToken(t, refreshToken)
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
	user.ID = userID
	test.AuthorID = userID

	accessToken, refreshToken, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
	}
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
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
			name: "Success: GetRefreshToken test",
			args: args{
				id:    foundTest.ID,
				token: accessToken,
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
			name: "Error: with invalid accessToken",
			args: args{
				id:    foundTest.ID,
				token: "bad accessToken",
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "No test found",
			args: args{
				id:    0,
				token: accessToken,
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

			r := gin.Default()
			r.Use(sessions.Sessions("session", mockHandlers.store))
			r.GET("/api/v1/tests/:id", setSessionMiddleware(t, tt.args.token), mockHandlers.authMiddleware, mockHandlers.GetTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				trueStatus := w.Code == tt.want.status

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
		helperDeleteRefreshTokenByToken(t, refreshToken)
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

	accessToken, refreshToken, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
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
			name: "Success: GetRefreshToken all tests with default pagination",
			args: args{
				repo:             mockRepo,
				createByQuantity: 10,
				params: domain.GetAllTestsRequest{
					PageID:   1,
					PageSize: 10,
				},
				token: accessToken,
			},
			want: want{
				createdByQuantity: 10,
				status:            http.StatusOK,
			},
		},
		{
			name: "Success: GetRefreshToken 1 page of tests but in db there are more than 10 tests",
			args: args{
				repo:             mockRepo,
				createByQuantity: 12,
				params: domain.GetAllTestsRequest{
					PageID:   1,
					PageSize: 10,
				},
				token: accessToken,
			},
			want: want{
				createdByQuantity: 10,
				status:            http.StatusOK,
			},
		},
		{
			name: "Success: GetRefreshToken 3 page with 2 tests, in db there are 22 tests",
			args: args{
				repo:             mockRepo,
				createByQuantity: 22,
				params: domain.GetAllTestsRequest{
					PageID:   3,
					PageSize: 10,
				},
				token: accessToken,
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
				token: accessToken,
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

			q := req.URL.Query()
			q.Add("page_id", strconv.Itoa(tt.args.params.PageID))
			q.Add("page_size", strconv.Itoa(tt.args.params.PageSize))
			req.URL.RawQuery = q.Encode()

			r := gin.Default()
			r.Use(sessions.Sessions("session", mockHandlers.store))
			r.GET("/api/v1/tests", setSessionMiddleware(t, tt.args.token), mockHandlers.authMiddleware, mockHandlers.GetAllTestsByUserID)

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
		helperDeleteRefreshTokenByToken(t, refreshToken)
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

	accessToken, refreshToken, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
	}
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
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
				token: accessToken,
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
				token: accessToken,
				input: badJSON,
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Error: invalid accessToken",
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

			r := gin.Default()
			r.Use(sessions.Sessions("session", mockHandlers.store))
			r.PUT("/api/v1/tests/:id", setSessionMiddleware(t, tt.args.token), mockHandlers.authMiddleware, mockHandlers.UpdateTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.want.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteTestByID(t, testIDZero)
		helperDeleteTestByID(t, testID)
		helperDeleteUserByID(t, userID)
		helperDeleteRefreshTokenByToken(t, refreshToken)
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

	accessToken, refreshToken, err := mockServices.Auth.SignIn(ctx, user)
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
	}
	if err != nil {
		t.Fatalf("error generating accessToken: %v", err)
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
				token: accessToken,
				id:    testID2,
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "Fail: Delete test that does not exist",
			args: args{
				token: accessToken,
				id:    2134234234,
			},
			want: want{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Error: invalid accessToken",
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
			r.Use(sessions.Sessions("session", mockHandlers.store))
			r.DELETE("/api/v1/tests/:id", setSessionMiddleware(t, tt.args.token), mockHandlers.authMiddleware, mockHandlers.DeleteTestByID)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.want.status
			})
		})
	}

	t.Cleanup(func() {
		helperDeleteUserByID(t, userID)
		helperDeleteTestByID(t, testID1)
		helperDeleteTestByID(t, testID2)
		helperDeleteRefreshTokenByToken(t, refreshToken)
	})
}

func setSessionMiddleware(t *testing.T, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(accessTokenName, token)
		if err := session.Save(); err != nil {
			t.Fatalf("error saving session: %v", err)
		}
	}
}
