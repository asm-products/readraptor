package readraptor

import (
	"os"

	"github.com/codegangsta/martini"
	"github.com/cupcake/gokiq"
	"github.com/garyburd/redigo/redis"
	workers "github.com/jrallison/go-workers"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func setupMartini(root string) *martini.Martini {
	m := martini.New()

	// database
	InitDb(os.Getenv("DATABASE_URL"))

	// Sessions Cookie store
	store := sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30,
	})
	m.Use(sessions.Sessions("rr_session", store))
	m.Use(sessionauth.SessionUser(GuestAccount))
	sessionauth.RedirectUrl = "/login"
	sessionauth.RedirectParam = "return"

	// middleware
	m.Use(ReqLogger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public", martini.StaticOptions{
		Prefix:      "assets",
		SkipLogging: true,
	}))
	m.Use(render.Renderer())

	// routes
	r := martini.NewRouter()
	r.Get("/", RedirectAuthenticated("/account"), func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	r.Get("/signout", sessionauth.LoginRequired, GetSignout)

	r.Post("/accounts", PostAccounts)
	r.Get("/account", sessionauth.LoginRequired, GetAccount)
	r.Post("/account/billing", sessionauth.LoginRequired, PostAccountBilling)

	r.Get("/setup", sessionauth.LoginRequired, GetSetup)

	r.Get("/confirm/:confirmation_token", GetConfirmAccount)
	r.Get("/t/:public_key/:article_id/:user_id/:signature.gif", GetTrackReadReceipts(root))
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
	gokiq.Client.Register(&NewAccountEmailJob{}, "default", 5)

	m.Map(gokiq.Client)

	m.Action(r.Handle)

	return m
}

func RunWeb(root string) {
	m := setupMartini(root)
	m.Run()
}
