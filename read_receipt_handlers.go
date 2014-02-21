package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

func GetTrackReadReceipts(dbmap *gorp.DbMap, params martini.Params, w http.ResponseWriter, r *http.Request) {
	if ensureSignatureMatch(dbmap, params, w, r) {
		account, err := FindAccountByUsername(dbmap, params["username"])
		if err != nil {
			panic(err)
		}

		err = TrackReadReceipt(dbmap, account, params["content_item_id"], params["user_id"])
		if err != nil {
			panic(err)
		}

		http.ServeFile(w, r, "public/tracking.gif")
	}
}

func ensureSignatureMatch(dbmap *gorp.DbMap, params martini.Params, w http.ResponseWriter, r *http.Request) bool {
	var apiKey string
	err := dbmap.SelectOne(&apiKey, `SELECT api_key FROM accounts WHERE username = $1`, params["username"])
	if err != nil {
		panic(err)
	}

	if params["signature"] != signature(apiKey, params["username"], params["content_item_id"], params["user_id"]) {
		http.NotFound(w, r)
		return false
	}
	return true
}

func signature(key, username, contentId, userId string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key + username + contentId + userId))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
