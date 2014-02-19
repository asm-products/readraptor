package main

type Account struct {
    Id       int64  `db:"id"`
    Created  int64  `db:"created_at"`
    Username string `db:"username"`
    ApiKey   string `db:"api_key"`
}
