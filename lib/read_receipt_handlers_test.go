package readraptor_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/asm-products/readraptor/fake"

	rr "github.com/asm-products/readraptor/lib"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

func Test_Tracking(t *testing.T) {
	dbmap := initTestDb(t)
	defer dbmap.Db.Close()

	account := MustCreateAccount(dbmap, "weasley@example.com")

	params := martini.Params{
		"public_key": account.PublicKey,
		"article_id": "article_1",
		"user_id":    "user_1",
		"signature":  rr.Signature(account.PrivateKey, account.PublicKey, "article_1", "user_1"),
	}

	req := &http.Request{}
	req.URL, _ = url.Parse("/t")
	rw := fake.New(t)
	rr.GetTrackReadReceipts("..")(params, rw, req)

	gif, _ := ioutil.ReadFile("../public/tracking.gif")
	rw.Assert(200, gif)
}
