package main

import (
	"context"
	"database/sql"
	"fmt"
	"linkr/internal/http/handlers"
	"linkr/internal/migrations"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	db := migrate()
	defer db.Close()
	addr, ok := os.LookupEnv("HTTP_ADDR")
	if !ok {
		addr = ":8080"
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: routes(db),
	}
	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("server stopped")
}

func routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/health", handlers.Health)
	r.Route("/{alias}", func(r chi.Router) {
		r.Get("/", handlers.Alias(db))
	})
	return r
}

func migrate() *sql.DB {
	dbName, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Println("Missing DB_USER")
		os.Exit(1)
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Println("Missing DB_PASSWORD")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:32768/postgres?sslmode=disable", dbName, dbPassword))
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	if err := migrations.Up(db); err != nil {
		panic(err)
	}
	return db
}
