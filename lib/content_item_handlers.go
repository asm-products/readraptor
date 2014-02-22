package readraptor

import (
	"encoding/json"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	"github.com/whatupdave/gokiq"
)

type ContentItemParams struct {
	Key       string           `json:"key"`
	Expected  []string         `json:"expected"`
	Callbacks []CallbackParams `json:"callbacks"`
}

type CallbackParams struct {
	Seconds int64  `json:"seconds"`
	Url     string `json:"url"`
}

func GetContentItems(dbmap *gorp.DbMap, params martini.Params) (string, int) {
	var ci ContentItem
	err := dbmap.SelectOne(&ci, "select * from content_items where key = $1", params["content_item_id"])
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

func PostContentItems(dbmap *gorp.DbMap, client *gokiq.ClientConfig, req *http.Request, account *Account) (string, int) {
	decoder := json.NewDecoder(req.Body)
	var p ContentItemParams
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	cid, err := InsertContentItem(dbmap, account.Id, p.Key)
	if err != nil {
		panic(err)
	}

	rids, err := AddReaders(dbmap, account.Id, cid, p.Expected)
	ScheduleCallbacks(client, rids, p.Callbacks)

	ci, err := FindContentItemWithReadReceipts(dbmap, cid)

	json, err := json.Marshal(map[string]interface{}{
		"content_item": ci,
	})
	if err != nil {
		panic(err)
	}
	return string(json), http.StatusCreated
}
