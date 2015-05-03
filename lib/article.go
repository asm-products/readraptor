package readraptor

import (
	"time"

	"github.com/coopernurse/gorp"
)

type Article struct {
	Id        int64     `db:"id"               json:"id"`
	AccountId int64     `db:"account_id"       json:"-"`
	Created   time.Time `db:"created_at"       json:"created"`
	Updated   time.Time `db:"updated_at"       json:"updated"`
	Key       string    `db:"key"              json:"key"`

	// TODO these fields shouldn't be in this struct, they're not in the articles
	// table
	FirstReadAt Timestamp `db:"first_read_at" json:"first_read_at,omitempty"`
	LastReadAt  Timestamp `db:"last_read_at"  json:"last_read_at,omitempty"`

	Delivered []string `json:"delivered,omitempty"`
	Pending   []string `json:"pending,omitempty"`
}

// Returns 0 if not found
func FindArticleIdByKey(dbmap *gorp.DbMap, accountId int64, key string) (id int64, err error) {
	return dbmap.SelectInt(`select id from articles where account_id = $1 and key = $2`, accountId, key)
}

func FindArticleWithReadReceipts(dbmap *gorp.DbMap, id int64) (*Article, error) {
	var a Article
	err := dbmap.SelectOne(&a, "select * from articles where id = $1", id)
	a.AddReadReceipts(dbmap)

	return &a, err
}

func (c *Article) AddReadReceipts(dbmap *gorp.DbMap) {
	var delivered []string
	_, err := dbmap.Select(&delivered, `
        select readers.distinct_id
        from articles
           inner join read_receipts on read_receipts.article_id = articles.id
           inner join readers on read_receipts.reader_id = readers.id
        where articles.id = $1
				  and last_read_at > articles.updated_at`, c.Id)
	if err != nil {
		panic(err)
	}
	c.Delivered = delivered

	var pending []string
	_, err = dbmap.Select(&pending, `
        select readers.distinct_id
        from articles
           inner join read_receipts on read_receipts.article_id = articles.id
           inner join readers on read_receipts.reader_id = readers.id
        where articles.id = $1
				  and (last_read_at is null or
						   last_read_at < articles.updated_at)`, c.Id)
	if err != nil {
		panic(err)
	}
	c.Pending = pending
}

func AddArticleReaders(dbmap *gorp.DbMap, accountId, articleId int64, expected []string) (rids []int64, err error) {
	for _, expectedReader := range expected {
		var rid int64
		rid, err = InsertReader(dbmap, accountId, expectedReader)
		if err != nil {
			return
		}
		rids = append(rids, rid)

		_, err = InsertExpectedReader(dbmap, articleId, rid)
		if err != nil {
			return
		}
	}
	return
}

func UpsertArticle(dbmap *gorp.DbMap, accountId int64, key string) (int64, error) {
	id, err := dbmap.SelectNullInt(`
        with s as (
            update articles set updated_at = $3 where account_id = $1 and key = $2 returning id
        ), i as (
            insert into articles ("account_id", "key", "created_at", "updated_at")
            select $1, $2, $3, $3
            where not exists (select 1 from s)
            returning id
        )
        select id from i union all select id from s;
    `, accountId,
		key,
		time.Now().UTC(),
	)
	if err != nil {
		return -1, err
	}

	iid, err := id.Value()

	return iid.(int64), err
}
