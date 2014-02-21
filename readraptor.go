package main

import (
	"os"

	"github.com/codegangsta/martini"
	workers "github.com/jrallison/go-workers"
)

var m *martini.Martini

func setupMartini() *martini.Martini {
	m := martini.New()

	// database
	db := InitDb(os.Getenv("DATABASE_URL"))

	// middleware
	m.Use(ReqLogger())
	m.Use(martini.Recovery())

	// routes
	r := martini.NewRouter()
	r.Post("/accounts", PostAccounts)
	r.Get("/t/:username/:content_item_id/:user_id/:signature.gif", GetTrackReadReceipts)
	r.Get("/content_items/:content_item_id", AuthAccount, GetContentItems)
	r.Post("/content_items", AuthAccount, PostContentItems)

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

func main() {
	m := setupMartini()
	m.Run()
}
