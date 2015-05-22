package readraptor

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

type ReadReceiptParams struct {
	Key        string `json:"key"`
	PublicKey  string `json:"public_key"`
	DistinctId string `json:"distinct_id"`
	Signature  string `json:"signature"`
}

/**
 * @api {get} /t/:username/:article_id/:user_id/:signature.gif Read article
 * @apiName GetTrack
 * @apiGroup Reading
 *
 */
func GetTrackReadReceipts(root string) func(params martini.Params, w http.ResponseWriter, r *http.Request) {
	return func(params martini.Params, w http.ResponseWriter, r *http.Request) {
		trackRead(params["public_key"], params["article_id"], params["user_id"], params["signature"], w, r)

		f, err := os.Open(root + "/public/tracking.gif")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		http.ServeContent(w, r, "tracking.gif", time.Time{}, f)
	}
}

func PostReadReceipts(params martini.Params, w http.ResponseWriter, r *http.Request) (string, int) {
	decoder := json.NewDecoder(r.Body)
	var p ReadReceiptParams
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	trackRead(p.PublicKey, p.Key, p.DistinctId, p.Signature, w, r)

	return "", http.StatusOK
}

func trackRead(publicKey, articleId, distinctId, signature string, w http.ResponseWriter, r *http.Request) {
	if ensureSignatureMatch(publicKey, articleId, distinctId, signature, w, r) {
		account, err := FindAccountByPublicKey(publicKey)
		if err != nil {
			panic(err)
		}

		err = TrackReadReceipt(dbmap, account, articleId, distinctId)
		if err != nil {
			panic(err)
		}
	}

}

func ensureSignatureMatch(publicKey, articleId, distinctId, signature string, w http.ResponseWriter, r *http.Request) bool {
	var privateKey string
	err := dbmap.SelectOne(&privateKey, `SELECT private_key FROM accounts WHERE public_key = $1`, publicKey)
	if err != nil {
		panic(err)
	}

	if signature != Signature(privateKey, publicKey, articleId, distinctId) {
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
