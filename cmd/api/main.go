package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"linkr/internal/migrations"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	migrate()
	addr, ok := os.LookupEnv("HTTP_ADDR")
	if !ok {
		addr = ":8080"
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: routes(),
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

func routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	return r
}

func migrate() {
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
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	if err := migrations.Up(db); err != nil {
		panic(err)
	}
}
