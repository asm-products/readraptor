package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "os"
)

func openDb() *sql.DB {
    connection := os.Getenv("DATABASE_URL")
    db, err := sql.Open("postgres", connection)
    if err != nil {
        panic(err)
    }

    return db
}
