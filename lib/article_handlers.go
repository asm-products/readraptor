package readraptor

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/cupcake/gokiq"
	"github.com/lib/pq"
	"github.com/technoweenie/grohl"
)

type ArticleParams struct {
	Key        string           `json:"key"`
	Recipients []string         `json:"recipients"`
	Callbacks  []CallbackParams `json:"via"`
}

type CallbackParams struct {
	At         int64    `json:"at"`
	Recipients []string `json:"recipients"`
	Url        string   `json:"url"`
}

func GetArticles(params martini.Params) (string, int) {
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

func PostArticles(client *gokiq.ClientConfig, req *http.Request, account *Account) (string, int) {
	decoder := json.NewDecoder(req.Body)
	var p ArticleParams
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	cid, err := InsertArticle(dbmap, account.Id, p.Key)
	if _, ok := err.(*pq.Error); ok {
		if strings.Index(err.Error(), `duplicate key value violates unique constraint "articles_key_key"`) == -1 {
			panic(err)
		}
	}

	grohl.Log(grohl.Data{
		"account":  account.Id,
		"register": p.Key,
		"readers":  p.Recipients,
	})

	rids, err := AddArticleReaders(dbmap, account.Id, cid, p.Recipients)
	for _, callback := range p.Callbacks {
		at := time.Unix(callback.At, 0)

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
