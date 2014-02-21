package main

import (
	"time"

	"github.com/coopernurse/gorp"
)

type ContentItem struct {
	Id        int64     `db:"id"`
	AccountId int64     `db:"account_id"`
	Created   time.Time `db:"created_at"`
	Key       string    `db:"key"`
}

func InsertContentItem(dbmap *gorp.DbMap, accountId int64, key string) (int64, error) {
	id, err := dbmap.SelectNullInt(`
        with s as (
            select id from content_items where account_id = $1 and key = $2
        ), i as (
            insert into content_items ("account_id", "key", "created_at")
            select $1, $2, $3
            where not exists (select 1 from s)
            returning id
        )
        select id from i union all select id from s;
    `, accountId,
		key,
		time.Now(),
	)
	if err != nil {
		return -1, err
	}

	iid, err := id.Value()

	return iid.(int64), err
}
