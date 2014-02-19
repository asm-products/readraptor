package main

import (
    "crypto/sha1"
    "database/sql"
    "fmt"
    "github.com/codegangsta/martini"
    _ "github.com/lib/pq"
    "github.com/technoweenie/grohl"
    "net/http"
)

func GetTracking(db *sql.DB, params martini.Params, w http.ResponseWriter, r *http.Request) {
    grohl.Log(grohl.Data{
        "username": params["username"],
        "content":  params["content_id"],
        "user":     params["user_id"],
        "sig":      params["signature"],
    })

    var apiKey string
    err := db.QueryRow(`SELECT api_key FROM accounts WHERE username = $1`, params["username"]).Scan(&apiKey)
    if err != nil {
        panic(err)
    }

    if params["signature"] != signature(apiKey, params["username"], params["content_id"], params["user_id"]) {
        http.NotFound(w, r)
    } else {
        http.ServeFile(w, r, "public/tracking.gif")
    }
}

func signature(key, username, contentId, userId string) string {
    hasher := sha1.New()
    hasher.Write([]byte(key + username + contentId + userId))
    return fmt.Sprintf("%x", hasher.Sum(nil))
}
