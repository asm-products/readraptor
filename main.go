package main

import (
    "encoding/json"
    "github.com/codegangsta/martini"
    "github.com/coopernurse/gorp"
    "os"
    // "github.com/martini-contrib/auth"
    "github.com/jrallison/go-workers"
    "net/http"
)

var m *martini.Martini

func setupMartini() *martini.Martini {
    m := martini.New()

    // database
    db := initDb(os.Getenv("DATABASE_URL"))

    // middleware
    m.Use(ReqLogger())
    m.Use(martini.Recovery())

    // routes
    r := martini.NewRouter()
    r.Post("/accounts", PostAccounts)
    r.Get("/t/:username/:content_id/:user_id/:signature.gif", GetTracking)

    // go-workers stats
    workers.Configure(map[string]string{
        "process": "web",
        "server":  "localhost:6379",
    })
    r.Get("/workers/stats", workers.Stats)

    // Inject database
    m.Map(db)

    m.Action(r.Handle)

    return m
}

func PostAccounts(dbmap *gorp.DbMap, req *http.Request) (string, int) {
    if err := req.ParseForm(); err != nil {
        panic(err)
    }

    account := NewAccount(req.PostForm["username"][0])
    err := dbmap.Insert(account)
    if err != nil {
        panic(err)
    }

    json, err := json.Marshal(account)
    if err != nil {
        panic(err)
    }
    return string(json), 201
}

func main() {
    m := setupMartini()
    m.Run()
}
