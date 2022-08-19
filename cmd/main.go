package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/server"
	"github.com/popeskul/qna-go/internal/services"
	"github.com/popeskul/qna-go/internal/transport/rest"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	ConfigDir  = "configs"
	ConfigFile = "config"
)

// @title Qna API
// @version 1.0
// @description Qna API
// @host localhost:8080
// @basePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg, err := config.New(ConfigDir, ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

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

	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	handlers := rest.NewHandler(service)

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
