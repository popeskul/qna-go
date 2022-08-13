package postgres

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/popeskul/qna-go/internal/db"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func TestMain(m *testing.M) {
	if err := changeDirToRoot(); err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestNewPostgresConnection(t *testing.T) {
	type args struct {
		cfg db.ConfigDB
	}
	tests := []struct {
		name string
		args args
		want *sql.DB
	}{
		{
			name: "Success",
			args: args{
				cfg: db.ConfigDB{
					Host:     "localhost",
					Port:     5432,
					User:     "postgres",
					Password: os.Getenv("DB_PASSWORD"),
					DBName:   "postgres",
					SSLMode:  "disable",
				},
			},
			want: &sql.DB{},
		},
		{
			name: "Fail",
			args: args{
				cfg: db.ConfigDB{
					Host:     "localhost",
					Port:     5432,
					User:     "postgres",
					Password: "postgres",
					DBName:   "qna",
					SSLMode:  "disable",
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewPostgresConnection(tt.args.cfg)

			if reflect.TypeOf(db) != reflect.TypeOf(tt.want) {
				t.Errorf("NewPostgresConnection() = %v, want %v", db, tt.want)
			}

			if err == nil && reflect.TypeOf(db) == reflect.TypeOf(tt.want) {
				err = db.Ping()

				if err != nil {
					t.Errorf("NewPostgresConnection() = %v, want %v", db, tt.want)
				}
			}
		})
	}
}

func changeDirToRoot() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "./../../../")
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
