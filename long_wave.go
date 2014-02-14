package main

import (
    "github.com/codegangsta/martini"
    "github.com/technoweenie/grohl"
    "net/http"
)

func NewMartini() *martini.ClassicMartini {
    r := martini.NewRouter()
    basic := martini.New()
    basic.Use(ReqLogger())
    basic.Use(martini.Recovery())
    basic.Action(r.Handle)
    return &martini.ClassicMartini{basic, r}
}

func GetTracker(params martini.Params, w http.ResponseWriter, r *http.Request) {
    grohl.Log(grohl.Data{
        "account": params["account_id"],
        "content": params["content_id"],
        "user":    params["user_id"],
    })

    http.ServeFile(w, r, "public/tracking.gif")
}

func main() {
    m := NewMartini()
    m.Get("/:account_id/:content_id/:user_id.gif", GetTracker)
    m.Run()
}
