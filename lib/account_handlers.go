package readraptor

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/cupcake/gokiq"
	pq "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"github.com/technoweenie/grohl"
)

func GetAccount(r render.Render, user sessionauth.User) {
	data := struct {
		StripeKey string
		Account   *Account
	}{
		os.Getenv("STRIPE_PUBLISHABLE"),
		user.(*Account),
	}
	r.HTML(200, "account", data)
}

func PostAccounts(client *gokiq.ClientConfig, req *http.Request) (string, int) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	account := NewAccount(req.PostForm["email"][0])
	err := dbmap.Insert(account)
	if err != nil {
		if _, ok := err.(*pq.Error); ok {
			if strings.Index(err.Error(), `duplicate key value violates unique constraint "accounts_email_key"`) > -1 {
				return "Email is already taken", http.StatusBadRequest
			}
		}
		panic(err)
	}

	json, err := json.Marshal(map[string]interface{}{
		"account": account,
	})
	if err != nil {
		panic(err)
	}

	account.SendNewAccountEmail(client)

	return string(json), http.StatusCreated
}

func PostAccountBilling(rw http.ResponseWriter, req *http.Request, user sessionauth.User) {
	account := user.(*Account)

	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	token := req.PostForm["stripeToken"][0]

	err := account.CreateStripeCustomer(token)
	if err != nil {
		panic(err)
	}

	http.Redirect(rw, req, "/setup", http.StatusFound)
}

func GetConfirmAccount(session sessions.Session, rw http.ResponseWriter, req *http.Request, params martini.Params) {
	account, err := FindAccountByConfirmationToken(params["confirmation_token"])
	if err != nil {
		panic(err)
	}

	_, err = dbmap.Exec(
		"update accounts set confirmed_at = $1 where id = $2",
		time.Now(),
		account.Id,
	)
	if err != nil {
		panic(err)
	}

	err = sessionauth.AuthenticateSession(session, account)
	if err != nil {
		panic(err)
	}

	http.Redirect(rw, req, "/account", http.StatusFound)
}

func GetSetup(r render.Render, user sessionauth.User, rw http.ResponseWriter, req *http.Request) {
	account := user.(*Account)
	if account.CustomerId == nil {
		http.Redirect(rw, req, "/account", http.StatusFound)
	} else {
		data := struct {
			Account *Account
		}{
			user.(*Account),
		}
		r.HTML(200, "setup", data)
	}
}

func AuthAccount(rw http.ResponseWriter, req *http.Request, c martini.Context) {
	splits := strings.Split(req.Header.Get("Authorization"), " ")
	dec, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		panic(err)
	}

	apiKey := string(dec[:len(dec)-1])

	var account Account
	err = dbmap.SelectOne(&account, "select * from accounts where private_key = $1 limit 1", apiKey)
	if err == sql.ErrNoRows {
		rw.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
		http.Error(rw, "Not Authorized", http.StatusUnauthorized)
	} else if err != nil {
		panic(err)
	}

	c.Map(&account)

	grohl.AddContext("account", account.Id)
}

func GetSignout(r render.Render, user sessionauth.User, session sessions.Session) {
	sessionauth.Logout(session, user)
	r.Redirect("/", http.StatusFound)
}

func RedirectAuthenticated(path string) martini.Handler {
	return func(r render.Render, user sessionauth.User, req *http.Request) {
		if user.IsAuthenticated() {
			r.Redirect(path, http.StatusFound)
		}
	}
}
