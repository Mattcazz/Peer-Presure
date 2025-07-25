package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mattcazz/Peer-Presure.git/db"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("HOST")
	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PASSWORD")
	dbName := os.Getenv("DATABASE")
	dbPort := os.Getenv("PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName) //postgres://USER:PASSWORD@HOST:PORT/DATABASE?OPTIONS

	db, err := db.New(dsn, 30, 30, "15m")

	if err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
