package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/popeskul/qna-go/internal/config"
	"github.com/popeskul/qna-go/internal/db"
	"github.com/popeskul/qna-go/internal/db/postgres"
	"github.com/popeskul/qna-go/internal/repository"
	"github.com/popeskul/qna-go/internal/server"
	"github.com/popeskul/qna-go/internal/services"
	"github.com/popeskul/qna-go/internal/transport/rest"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ConfigDir  = "configs"
	ConfigFile = "config"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg, err := config.New(ConfigDir, ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

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
