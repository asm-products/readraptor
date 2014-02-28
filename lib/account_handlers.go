package readraptor

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/codegangsta/martini"
	"github.com/cupcake/gokiq"
	"github.com/lib/pq"
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

func PostAccounts(client *gokiq.ClientConfig, rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	email := strings.Replace(req.PostForm["email"][0], " ", "", -1)

	if strings.Index(email, "@") == -1 {
		rw.Write([]byte("Email is invalid"))
		return
	}

	account := NewAccount(email)
	err := dbmap.Insert(account)
	if err != nil {
		if _, ok := err.(*pq.Error); ok {
			if strings.Index(err.Error(), `duplicate key value violates unique constraint "accounts_email_key"`) > -1 {
				rw.Write([]byte("Email is already taken"))
				return
			}
		}
		panic(err)
	}

	account.SendNewAccountEmail(client)

	template, err := template.ParseFiles("templates/welcome.txt.tmpl")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = template.Execute(&buf, nil)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(rw, buf.String())
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
			Account   *Account
			Signature string
		}{
			user.(*Account),
			signature(account.PrivateKey, account.PublicKey, "article_1", "user_1"),
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

	privateKey := string(dec[:len(dec)-1])

	var account Account
	err = dbmap.SelectOne(&account, "select * from accounts where private_key = $1 limit 1", privateKey)
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
