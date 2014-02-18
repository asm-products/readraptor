package main

import (
    "crypto/sha1"
    "fmt"
    "github.com/codegangsta/martini"
    "github.com/martini-contrib/auth"
    "github.com/technoweenie/grohl"
    "net/http"
)

func NewMartini() *martini.ClassicMartini {
    r := martini.NewRouter()
    m := martini.New()
    m.Use(ReqLogger())
    m.Use(martini.Recovery())
    m.Action(r.Handle)
    return &martini.ClassicMartini{m, r}
}

func TrackHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
    sig := signature(params["username"], params["content_id"], params["user_id"])

    grohl.Log(grohl.Data{
        "username": params["username"],
        "content":  params["content_id"],
        "user":     params["user_id"],
        "sig":      sig,
    })

    http.ServeFile(w, r, "public/tracking.gif")
}

func signature(username, contentId, userId string) string {
    key := "api_3c12d9556813"
    hasher := sha1.New()
    hasher.Write([]byte(key + username + contentId + userId))
    return fmt.Sprintf("%x", hasher.Sum(nil))
}

func main() {
    m := NewMartini()
    m.Get("/t/:username/:content_id/:user_id/:checksum.gif", auth.Basic("whatupdave", ""), TrackHandler)
    m.Run()
}
