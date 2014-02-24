package readraptor

import (
	"os"

	"github.com/codegangsta/martini"
	"github.com/cupcake/gokiq"
	"github.com/garyburd/redigo/redis"
	workers "github.com/jrallison/go-workers"
	"github.com/martini-contrib/render"
)

func setupMartini(root string) *martini.Martini {
	m := martini.New()

	// database
	InitDb(os.Getenv("DATABASE_URL"))

	// middleware
	m.Use(ReqLogger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public", martini.StaticOptions{
		Prefix:      "public",
		SkipLogging: true,
	}))
	m.Use(render.Renderer())

	// routes
	r := martini.NewRouter()
	r.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})
	r.Post("/accounts", PostAccounts)
	r.Get("/t/:username/:article_id/:user_id/:signature.gif", GetTrackReadReceipts(root))
	r.Get("/articles/:article_id", AuthAccount, GetArticles)
	r.Post("/articles", AuthAccount, PostArticles)

	// go-workers stats
	workers.Configure(map[string]string{
		"process": "web",
		"server":  "localhost:6379",
	})
	r.Get("/workers/stats", workers.Stats)

	// Inject database
	m.Map(dbmap)

	// Inject gokiq client
	gokiq.Client.RedisNamespace = "rr"
	gokiq.Workers.RedisPool = redis.NewPool(RedisConnect(os.Getenv("REDIS_URL")), 1)
	gokiq.Client.Register(&UserCallbackJob{}, "default", 5)

	m.Map(gokiq.Client)

	m.Action(r.Handle)

	return m
}

func RunWeb(root string) {
	m := setupMartini(root)
	m.Run()
}
