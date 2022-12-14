package v1

import (
	"context"
	"database/sql"
	sessionsPostgres "github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/popeskul/cache"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/hash"
	"github.com/popeskul/qna-go/internal/logger"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/repository/sessions"
	"github.com/popeskul/qna-go/internal/services"
	"github.com/popeskul/qna-go/internal/token"
	"github.com/popeskul/qna-go/internal/util"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	mockDB       *sql.DB
	mockRepo     *repository.Repository
	mockHandlers *Handlers
	mockServices *services.Service
	cfg          *config.Config
)

func TestMain(m *testing.M) {
	if err := util.ChangeDir("../../"); err != nil {
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

	pasetoMaker, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		log.Fatal(err)
	}

	hashManager, err := hash.NewHash(cfg.HashSalt)
	if err != nil {
		log.Fatal(err)
	}

	store, err := sessionsPostgres.NewStore(db, []byte(cfg.Session.Secret))
	if err != nil {
		log.Fatal(err)
	}

	d, err := time.ParseDuration(cfg.Cache.TTL)
	if err != nil {
		log.Fatal(err)
	}
	cache := cache.New(d)

	sessionManager := sessions.NewRepoSessions(db)

	mockRepo = repository.NewRepository(db)
	mockServices = services.NewService(mockRepo, pasetoMaker, hashManager, cache, sessionManager)
	mockHandlers = NewHandler(mockServices, store, logger.GetLogger())

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func randomUser() domain.User {
	return domain.User{
		Name:     util.RandomString(10),
		Email:    util.RandomString(10) + "@gmail.com",
		Password: util.RandomString(10),
	}
}

func randomTest() domain.Test {
	return domain.Test{
		Title:    util.RandomString(10),
		AuthorID: int(util.RandomInt(1, 100)),
	}
}

func helperCreatUser(t *testing.T, ctx context.Context, user domain.User) {
	err := mockServices.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}
}

func helperCreateTest(t *testing.T, userID int, test domain.Test) int {
	t.Helper()

	var id int
	if err := mockDB.QueryRow("INSERT INTO tests (title, author_id) VALUES ($1, $2) RETURNING id", test.Title, userID).Scan(&id); err != nil {
		t.Errorf("error inserting test: %v", err)
	}

	return id
}

func helperDeleteUserByID(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		t.Errorf("error deleting user: %v", err)
	}
}

func helperDeleteUserByEmail(t *testing.T, email string) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM users WHERE email = $1", email); err != nil {
		t.Errorf("error deleting user: %v", err)
	}
}

func helperDeleteTestByID(t *testing.T, id int) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM tests WHERE id = $1", id); err != nil {
		t.Errorf("error deleting test: %v", err)
	}
}

func helperDeleteRefreshTokenByToken(t *testing.T, token string) {
	t.Helper()
	if _, err := mockDB.Exec("DELETE FROM refresh_tokens WHERE token = $1", token); err != nil {
		t.Errorf("error deleting refresh token: %v", err)
	}
}

func helperFindTestByTitle(t *testing.T, title string) domain.Test {
	t.Helper()
	var test domain.Test
	if err := mockDB.QueryRow("SELECT id, title, author_id FROM tests WHERE title = $1", title).Scan(&test.ID, &test.Title, &test.AuthorID); err != nil {
		t.Errorf("error finding test: %v", err)
	}

	return test
}

func findUserIDByEmail(email string) (int, error) {
	var id int
	if err := mockDB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
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
	cfg.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	cfg.HashSalt = os.Getenv("HASH_SALT")
	cfg.Session.Secret = os.Getenv("SESSION_SECRET")

	return cfg, nil
}
