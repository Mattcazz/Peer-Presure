package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mattcazz/Peer-Presure.git/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName) //postgres://USER:PASSWORD@HOST:PORT/DATABASE?OPTIONS

	db, err := db.New(dsn, 30, 30, "15m")

	if err != nil {
		log.Fatal(err)
	}

	// -------------------- Migrations -------------------- //
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres", driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	err = m.Up();

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	// ---------------------------------------------------- //

	server := NewApiServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
