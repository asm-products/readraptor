package main

import (
    "crypto/sha1"
    "fmt"
    "os"
    "time"
)

type Account struct {
    Id       int64     `db:"id"`
    Created  time.Time `db:"created_at"`
    Username string    `db:"username"`
    ApiKey   string    `db:"api_key"`
}

func NewAccount(username string) *Account {
    account := &Account{
        Username: username,
        Created:  time.Now(),
    }
    account.ApiKey = genApiKey(username)

    return account
}

func genApiKey(username string) string {
    key := os.Getenv("API_GEN_SECRET")
    hasher := sha1.New()
    hasher.Write([]byte(key + username))
    return fmt.Sprintf("%x", hasher.Sum(nil))
}
