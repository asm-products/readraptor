package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

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

func AuthAccount(dbmap *gorp.DbMap, rw http.ResponseWriter, req *http.Request) {
	splits := strings.Split(req.Header.Get("Authorization"), " ")
	dec, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		panic(err)
	}

	apiKey := string(dec[:len(dec)-1])

	username, err := dbmap.SelectNullStr("select username from accounts where api_key = $1", apiKey)
	if err != nil {
		panic(err)
	}

	if !username.Valid {
		rw.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
		http.Error(rw, "Not Authorized", http.StatusUnauthorized)
	}
}
