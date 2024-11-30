package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

var DB *sql.DB

func InitDB() {
    connStr := "host=localhost port=5432 user=postgres password=password dbname=beauty_salon sslmode=disable"
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }

    log.Println("Connected to the database!")
}
