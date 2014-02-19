package main

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func openDb(connection string) *sql.DB {
    db, err := sql.Open("postgres", connection)
    if err != nil {
        panic(err)
    }

    return db
}
