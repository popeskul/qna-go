package v1

import (
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

var handlers *Handlers
var dbConn *sql.DB
var s *services.Service

func TestMain(m *testing.M) {
	if err := changeDirToRoot(); err != nil {
		log.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := newDBConnection(cfg)
	dbConn = db
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	s = services.NewService(repo)
	handlers = NewHandler(s)

	m.Run()
}

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
			if _, err := dbConn.Exec("INSERT INTO users (email, encrypted_password, name) VALUES ($1, $2, $3) RETURNING id", mockUniqEmail, mockPassword, "-"); err != nil {
				t.Errorf("Some error occured. Err: %s", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewReader(tt.user))
			req.Header.Set("Content-Type", "application/json")

			r := gin.Default()
			r.POST("/api/v1/auth/sign-up", handlers.SignUp)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.status
			})

			t.Cleanup(func() {
				if _, err := dbConn.Exec("DELETE FROM users WHERE email IN ($1, $2)", mockEmail, mockUniqEmail); err != nil {
					t.Errorf("Some error occured. Err: %s", err)
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
			if _, err := s.CreateUser(domain.SignUpInput{
				Email:             mockEmail,
				EncryptedPassword: mockPassword,
			}); err != nil {
				t.Errorf("Some error occured. Err: %s", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(tt.user))
			req.Header.Set("Content-Type", "application/json")

			r := gin.Default()
			r.POST("/api/v1/auth/sign-in", handlers.SignIn)

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				return w.Code == tt.status
			})

			t.Cleanup(func() {
				if _, err := dbConn.Exec("DELETE FROM users WHERE email IN ($1)", mockEmail); err != nil {
					t.Errorf("Some error occured. Err: %s", err)
				}
			})
		})
	}
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func newDBConnection(cfg *config.Config) (*sql.DB, error) {
	return postgres.NewPostgresConnection(db.ConfigDB{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
}

func loadConfig() (*config.Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg, err := config.New("configs", "config")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return cfg, nil
}

func changeDirToRoot() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "./../../../../")
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}
