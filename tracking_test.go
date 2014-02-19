package main

import (
    // "database/sql"
    "github.com/asm-products/long-wave/fake"
    "github.com/codegangsta/martini"
    "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
    "io/ioutil"
    "net/http"
    "net/url"
    "reflect"
    "testing"
    "time"
)

func Test_Tracking(t *testing.T) {
    // initialize the DbMap
    dbmap := initDb(t)
    defer dbmap.Db.Close()

    // delete any existing rows
    err := dbmap.TruncateTables()
    checkErr(t, err, "TruncateTables failed")

    err = dbmap.Insert(&Account{
        Username: "weasley",
        ApiKey:   "api1234",
        Created:  time.Now().UnixNano(),
    })
    checkErr(t, err, "Insert failed")

    params := martini.Params{
        "username":   "weasley",
        "content_id": "content_1",
        "user_id":    "user_1",
        "signature":  signature("api1234", "weasley", "content_1", "user_1"),
    }

    req := &http.Request{}
    req.URL, _ = url.Parse("/t")
    rw := fake.New(t)
    GetTracking(dbmap.Db, params, rw, req)

    gif, _ := ioutil.ReadFile("public/tracking.gif")
    rw.Assert(200, gif)
}

func expect(t *testing.T, a interface{}, b interface{}) {
    if a != b {
        t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
    }
}

func initDb(t *testing.T) *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    db := openDb("postgres://localhost/lw_test?sslmode=disable")

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Account{}, "accounts").SetKeys(true, "Id")

    err := dbmap.CreateTablesIfNotExists()
    checkErr(t, err, "Create tables failed")

    return dbmap
}

func checkErr(t *testing.T, err error, message string) {
    if err != nil {
        t.Fatalf("%s – %s", err, message)
    }
}
