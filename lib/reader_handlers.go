package readraptor

import (
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
)

func GetReader(account *Account, params martini.Params) (string, int) {
	var readerId int64
	readerId, err := dbmap.SelectInt(`
        select id
        from readers
        where account_id = $1
          and distinct_id = $2;`, account.Id, params["distinct_id"])
	if err != nil {
		panic(err)
	}

	articles, err := UnreadArticles(dbmap, readerId)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(articles)
	if err != nil {
		panic(err)
	}

	return string(json), http.StatusOK
}
