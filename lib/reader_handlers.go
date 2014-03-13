package readraptor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

func GetReader(params martini.Params) (string, int) {
	readerId, err := strconv.ParseInt(params["reader_id"], 10, 64)
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
