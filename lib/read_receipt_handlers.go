package readraptor

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

/**
 * @api {get} /t/:username/:article_id/:user_id/:signature Read article
 * @apiName GetTrack
 * @apiGroup Reading
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "firstname": "John",
 *       "lastname": "Doe"
 *     }
 */
func GetTrackReadReceipts(root string) func(params martini.Params, w http.ResponseWriter, r *http.Request) {
	return func(params martini.Params, w http.ResponseWriter, r *http.Request) {
		if ensureSignatureMatch(params, w, r) {
			account, err := FindAccountByPublicKey(params["public_key"])
			if err != nil {
				panic(err)
			}

			err = TrackReadReceipt(dbmap, account, params["article_id"], params["user_id"])
			if err != nil {
				panic(err)
			}

			f, err := os.Open(root + "/public/tracking.gif")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			http.ServeContent(w, r, "tracking.gif", time.Time{}, f)
		}
	}
}

func ensureSignatureMatch(params martini.Params, w http.ResponseWriter, r *http.Request) bool {
	var privateKey string
	err := dbmap.SelectOne(&privateKey, `SELECT private_key FROM accounts WHERE public_key = $1`, params["public_key"])
	if err != nil {
		panic(err)
	}

	if params["signature"] != Signature(privateKey, params["public_key"], params["article_id"], params["user_id"]) {
		http.NotFound(w, r)
		return false
	}
	return true
}

func Signature(key, public_key, articleId, userId string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key + public_key + articleId + userId))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
