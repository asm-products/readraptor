package main

import (
    "crypto/sha1"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/codegangsta/martini"
    _ "github.com/lib/pq"
    "os"
    // "github.com/martini-contrib/auth"
    "net/http"
)

var m *martini.Martini

func setupMartini() *martini.Martini {
    m := martini.New()

    // database
    db := openDb(os.Getenv("DATABASE_URL"))

    // middleware
    m.Use(ReqLogger())
    m.Use(martini.Recovery())

    // routes
    r := martini.NewRouter()
    r.Post("/accounts", PostAccounts)
    r.Get("/t/:username/:content_id/:user_id/:signature.gif", GetTracking)

    // Inject database
    m.Map(db)

    m.Action(r.Handle)

    return m
}

func PostAccounts(db *sql.DB, req *http.Request) (string, int) {
    if err := req.ParseForm(); err != nil {
        panic(err)
    }

    username := req.PostForm["username"][0]
    apiKey := genApiKey(username)
    _, err := db.Query(
        `INSERT INTO accounts(username, api_key) VALUES ($1, $2)`,
        username,
        apiKey,
    )
    if err != nil {
        panic(err)
    }

    json, err := json.Marshal(&Account{
        Username: username,
        ApiKey:   apiKey,
    })
    if err != nil {
        panic(err)
    }
    return string(json), 201
}

func genApiKey(username string) string {
    key := os.Getenv("API_GEN_SECRET")
    hasher := sha1.New()
    hasher.Write([]byte(key + username))
    return fmt.Sprintf("%x", hasher.Sum(nil))
}

func main() {
    m := setupMartini()
    m.Run()
}
