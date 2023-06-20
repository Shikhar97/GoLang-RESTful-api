package database

import (
	"database/sql"
	"log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func New() *sql.DB {
    log.Println("connecting to db")
    database, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable");
    if err != nil {
        log.Fatalf("error connecting to db: %e", err)
    }
    return database
}
