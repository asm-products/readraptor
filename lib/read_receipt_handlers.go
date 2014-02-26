package readraptor

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

func GetTrackReadReceipts(root string) func(*gorp.DbMap, martini.Params, http.ResponseWriter, *http.Request) {
	return func(dbmap *gorp.DbMap, params martini.Params, w http.ResponseWriter, r *http.Request) {
		if ensureSignatureMatch(dbmap, params, w, r) {
			account, err := FindAccountByPublicKey(params["public_key"])
			if err != nil {
				panic(err)
			}

			err = TrackReadReceipt(dbmap, account, params["article_id"], params["user_id"])
			if err != nil {
				panic(err)
			}

			http.ServeFile(w, r, root+"/public/tracking.gif")
		}
	}
}

func ensureSignatureMatch(dbmap *gorp.DbMap, params martini.Params, w http.ResponseWriter, r *http.Request) bool {
	var apiKey string
	err := dbmap.SelectOne(&apiKey, `SELECT private_key FROM accounts WHERE public_key = $1`, params["public_key"])
	if err != nil {
		panic(err)
	}

	if params["signature"] != signature(apiKey, params["public_key"], params["article_id"], params["user_id"]) {
		http.NotFound(w, r)
		return false
	}
	return true
}

func signature(key, public_key, articleId, userId string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key + public_key + articleId + userId))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
