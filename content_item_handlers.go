package main

import (
	"encoding/json"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
)

type ContentItemParams struct {
	Key    string   `json:"key"`
	Unseen []string `json:"unseen"`
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

func PostContentItems(dbmap *gorp.DbMap, req *http.Request, account *Account) (string, int) {
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

	for _, expectedReader := range p.Unseen {
		rid, err := InsertReader(dbmap, account.Id, expectedReader)
		if err != nil {
			panic(err)
		}

		_, err = InsertExpectedReader(dbmap, cid, rid)
		if err != nil {
			panic(err)
		}
	}

	ci, err := FindContentItem(dbmap, cid)

	json, err := json.Marshal(map[string]interface{}{
		"content_item": ci,
	})
	if err != nil {
		panic(err)
	}
	return string(json), http.StatusCreated
}
