package v1

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"
)

var mockDB *sql.DB
var mockRepo *repository.Repository
var mockHandlers *Handlers
var mockServices *services.Service

func TestMain(m *testing.M) {
	if err := changeDirToRoot(); err != nil {
		log.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := newDBConnection(cfg)
	mockDB = db
	if err != nil {
		log.Fatalf("Some error occured. Err: %mockServices", err)
	}
	defer db.Close()

	mockRepo = repository.NewRepository(db)
	mockServices = services.NewService(mockRepo)
	mockHandlers = NewHandler(mockServices)

	os.Exit(m.Run())
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()

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

	cfg, err := config.New("configs", "test.config")
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
