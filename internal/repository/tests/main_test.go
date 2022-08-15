package tests

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

var mockDB *sql.DB
var mockRepo *RepositoryTests

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
		log.Fatal(err)
	}

	mockRepo = NewRepoTests(mockDB)

	os.Exit(m.Run())
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
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
	dir := path.Join(path.Dir(filename), "./../../../")
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}
