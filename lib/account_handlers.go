package readraptor

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	"github.com/cupcake/gokiq"
	pq "github.com/lib/pq"
	"github.com/technoweenie/grohl"
)

func PostAccounts(dbmap *gorp.DbMap, client *gokiq.ClientConfig, req *http.Request) (string, int) {
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
		fmt.Println(reflect.TypeOf(err), err.Error())
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

func GetConfirmAccount(dbmap *gorp.DbMap, rw http.ResponseWriter, req *http.Request, params martini.Params) {
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

	http.Redirect(rw, req, "/account", http.StatusFound)
}

func AuthAccount(dbmap *gorp.DbMap, rw http.ResponseWriter, req *http.Request, c martini.Context) {
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
