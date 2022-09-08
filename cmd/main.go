package main

import (
	"context"
	"fmt"
	"github.com/popeskul/qna-go/internal/repository/sessions"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sessionsPostgres "github.com/gin-contrib/sessions/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/hash"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/server"
	"github.com/popeskul/qna-go/internal/services"
	"github.com/popeskul/qna-go/internal/token"
	"github.com/popeskul/qna-go/internal/transport/rest"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	if err = runMigration(cfg); err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgresConnection(db.ConfigDB{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		DBName:   cfg.DB.DBName,
		Password: cfg.DB.Password,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()

	tokenManager, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
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

	sessionManager := sessions.NewRepoSessions(db)

	repo := repository.NewRepository(db)
	service := services.NewService(repo, tokenManager, hashManager, sessionManager)
	handlers := rest.NewHandler(service, store)

	srv := server.NewServer(&http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        handlers.Init(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})

	go func() {
		log.Fatal(srv.Run())
	}()

	log.Println("Starting server on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	fmt.Println("Server shutting down...")

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Failed to shutdown server: ", err)
	}

	if err = db.Close(); err != nil {
		log.Fatal("Failed to close database: ", err)
	}

	fmt.Println("Server stopped")
}

// initConfig reading from the .env file and environment variables and returns the config and error.
func initConfig() (*config.Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg, err := config.New("configs", "config")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	cfg.HashSalt = os.Getenv("HASH_SALT")
	cfg.Session.Secret = os.Getenv("SESSION_SECRET")

	return cfg, nil
}

// runMigration run the migration for the database.
func runMigration(cfg *config.Config) error {
	migrationPath := "file://schema"
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName, cfg.DB.SSLMode)
	m, err := migrate.New(migrationPath, dbConn)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil {
		// because of the way up and down works, we need to check for the ErrNoChange error
		// without this check, the application will panic if the database is already up-to-date
		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}
