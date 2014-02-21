package main

import (
	"encoding/json"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/technoweenie/grohl"
)

func GetContentItems(params martini.Params) (string, int) {
	grohl.Log(grohl.Data{
		"username": params["username"],
		"content":  params["content_item_id"],
		"user":     params["user_id"],
		"sig":      params["signature"],
	})

	json, err := json.Marshal(&ContentItem{})
	if err != nil {
		panic(err)
	}

	return string(json), http.StatusOK
}
