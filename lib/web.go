package readraptor

import (
	"os"

	"github.com/codegangsta/martini"
	workers "github.com/jrallison/go-workers"
	"github.com/whatupdave/gokiq"
)

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

	// Inject gokiq client
	gokiq.Client.RedisNamespace = "rr"
	gokiq.Client.Register(&UserCallbackJob{}, "default", 5)

	m.Map(gokiq.Client)

	m.Action(r.Handle)

	return m
}

func RunWeb() {
	m := setupMartini()
	m.Run()
}
