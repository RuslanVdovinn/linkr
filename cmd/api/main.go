package main

import (
	"context"
	"fmt"
	"linkr/internal/auth"
	"linkr/internal/http/handlers"
	"linkr/internal/http/middleware"
	"linkr/internal/migrations"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	token, err := auth.GenerateToken(1, "test@test.com", time.Hour)
	log.Println(err)
	log.Printf("Current token: %s", token)
	gormDB := connect()
	db, _ := gormDB.DB()
	defer db.Close()
	addr, ok := os.LookupEnv("HTTP_ADDR")
	if !ok {
		addr = ":8080"
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: routes(gormDB),
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

func routes(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))
	r.Use(middleware.AuthBearerMiddleware)
	r.Get("/health", handlers.Health)
	r.Route("/{alias}", func(r chi.Router) {
		r.Get("/", handlers.Alias(db))
	})
	r.Route("/links", func(r chi.Router) {
		r.Post("/", handlers.CreateLink(db))
		r.Get("/{alias}", handlers.GetAlias(db))
		r.Patch("/{alias}", handlers.PatchAlias(db))
		r.Delete("/{alias}", handlers.DeleteAlias(db))
	})
	return r
}

func connect() *gorm.DB {
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
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@localhost:32771/postgres?sslmode=disable", dbName, dbPassword)),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}
	if err := migrations.Up(sqlDB); err != nil {
		panic(err)
	}
	return db
}
