package main

import (
	"time"

	"github.com/coopernurse/gorp"
)

type Reader struct {
	Id         int64     `db:"id"`
	Created    time.Time `db:"created_at"`
	AccountId  int64     `db:"account_id"`
	DistinctId string    `db:"distinct_id"`
}

func InsertReader(dbmap *gorp.DbMap, accountId int64, distinctId string) (int64, error) {
	id, err := dbmap.SelectNullInt(`
        with s as (
            select id from readers where account_id = $1 and distinct_id = $2
        ), i as (
            insert into readers ("account_id", "distinct_id", "created_at")
            select $1, $2, $3
            where not exists (select 1 from s)
            returning id
        )
        select id from i union all select id from s;
    `, accountId,
		distinctId,
		time.Now(),
	)
	if err != nil {
		return -1, err
	}

	iid, err := id.Value()

	return iid.(int64), err
}
