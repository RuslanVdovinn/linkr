package main

import (
	"database/sql"
	"fmt"
	"linkr/internal/migrations"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
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
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:32768/linkr?sslmode=disable", dbName, dbPassword))
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
