package auth

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var mockDB *sql.DB
var mockRepo *RepositoryAuth

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

	mockRepo = NewRepoAuth(mockDB)

	os.Exit(m.Run())
}
