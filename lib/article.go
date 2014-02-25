package readraptor

import (
	"time"

	"github.com/coopernurse/gorp"
)

type Article struct {
	Id        int64     `db:"id"         json:"id"`
	AccountId int64     `db:"account_id" json:"-"`
	Created   time.Time `db:"created_at" json:"created"`
	Key       string    `db:"key"        json:"key"`

	Delivered []string `json:"delivered,omitempty"`
	Pending   []string `json:"undelivered,omitempty"`
}

func FindArticleWithReadReceipts(dbmap *gorp.DbMap, id int64) (*Article, error) {
	var a Article
	err := dbmap.SelectOne(&a, "select * from articles where id = $1", id)
	a.AddReadReceipts(dbmap)

	return &a, err
}

func (c *Article) AddReadReceipts(dbmap *gorp.DbMap) {
	var delivered []string
	selectReaders := `
        select readers.distinct_id
        from articles
           inner join read_receipts on read_receipts.article_id = articles.id
           inner join readers on read_receipts.reader_id = readers.id
        where articles.id = $1`

	_, err := dbmap.Select(&delivered, selectReaders, c.Id)
	if err != nil {
		panic(err)
	}
	c.Delivered = delivered

	var pending []string
	_, err = dbmap.Select(&pending, `
        select readers.distinct_id
        from articles
           inner join expected_readers on expected_readers.article_id = articles.id
           inner join readers on expected_readers.reader_id = readers.id
        where articles.id = $1
        except all `+selectReaders, c.Id)
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

func InsertArticle(dbmap *gorp.DbMap, accountId int64, key string) (int64, error) {
	id, err := dbmap.SelectNullInt(`
        with s as (
            select id from articles where account_id = $1 and key = $2
        ), i as (
            insert into articles ("account_id", "key", "created_at")
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
