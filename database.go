package main

import (
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
)

func initDb(connection string) *gorp.DbMap {
    db, err := sql.Open("postgres", connection)
    if err != nil {
        panic(err)
    }
    
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}    
    dbmap.AddTableWithName(Account{}, "accounts").SetKeys(true, "Id")
    
    return dbmap
}
