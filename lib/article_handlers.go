package readraptor

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	"github.com/cupcake/gokiq"
)

type ArticleParams struct {
	Key        string           `json:"key"`
	Recipients []string         `json:"recipients"`
	Callbacks  []CallbackParams `json:"callbacks"`
}

type CallbackParams struct {
	Delay      string   `json:"delay"`
	Recipients []string `json:"recipients"`
	Url        string   `json:"url"`
}

func GetArticles(dbmap *gorp.DbMap, params martini.Params) (string, int) {
	var ci Article
	err := dbmap.SelectOne(&ci, "select * from articles where key = $1", params["article_id"])
	ci.AddReadReceipts(dbmap)

	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(ci)
	if err != nil {
		panic(err)
	}

	return string(json), http.StatusOK
}

func PostArticles(dbmap *gorp.DbMap, client *gokiq.ClientConfig, req *http.Request, account *Account) (string, int) {
	decoder := json.NewDecoder(req.Body)
	var p ArticleParams
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	cid, err := InsertArticle(dbmap, account.Id, p.Key)
	if err != nil {
		panic(err)
	}

	rids, err := AddArticleReaders(dbmap, account.Id, cid, p.Recipients)
	for _, callback := range p.Callbacks {
		delay, err := time.ParseDuration(callback.Delay)
		if err != nil {
			panic(err)
		}
		at := time.Now().Add(delay)

		if callback.Recipients != nil {
			rids, err = AddArticleReaders(dbmap, account.Id, cid, callback.Recipients)
			if err != nil {
				panic(err)
			}
		}
		ScheduleCallbacks(client, rids, at, callback.Url)
	}

	ci, err := FindArticleWithReadReceipts(dbmap, cid)

	json, err := json.Marshal(map[string]interface{}{
		"article": ci,
	})
	if err != nil {
		panic(err)
	}
	return string(json), http.StatusCreated
}
