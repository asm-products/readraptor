package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
)

func PostAccounts(dbmap *gorp.DbMap, req *http.Request) (string, int) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	account := NewAccount(req.PostForm["username"][0])
	err := dbmap.Insert(account)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(account)
	if err != nil {
		panic(err)
	}
	return string(json), http.StatusCreated
}

func AuthAccount(dbmap *gorp.DbMap, rw http.ResponseWriter, req *http.Request, c martini.Context) {
	splits := strings.Split(req.Header.Get("Authorization"), " ")
	dec, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		panic(err)
	}

	apiKey := string(dec[:len(dec)-1])

	var account Account
	err = dbmap.SelectOne(&account, "select * from accounts where api_key = $1 limit 1", apiKey)
	if err == sql.ErrNoRows {
		rw.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
		http.Error(rw, "Not Authorized", http.StatusUnauthorized)
	} else if err != nil {
		panic(err)
	}

	c.Map(&account)
}
