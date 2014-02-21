package main

import (
	"time"

	"github.com/coopernurse/gorp"
)

type ExpectedReader struct {
    Reader
}

func InsertExpectedReader(dbmap *gorp.DbMap, contentId, readerId int64) (int64, error) {
	id, err := dbmap.SelectNullInt(`
        with s as (
            select id from expected_readers where content_item_id = $1 and reader_id = $2
        ), i as (
            insert into expected_readers ("content_item_id", "reader_id", "created_at")
            select $1, $2, $3
            where not exists (select 1 from s)
            returning id
        )
        select id from i union all select id from s;
    `, contentId,
		readerId,
		time.Now(),
	)
	if err != nil {
		return -1, err
	}

	iid, err := id.Value()

	return iid.(int64), err
}