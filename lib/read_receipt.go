package readraptor

import (
	"time"

	"github.com/coopernurse/gorp"
	"github.com/technoweenie/grohl"
)

type ReadReceipt struct {
	Id          int64      `db:"id"`
	Created     time.Time  `db:"created_at"`
	ArticleId   int64      `db:"article_id"`
	ReaderId    int64      `db:"reader_id"`
	FirstReadAt *time.Time `db:"first_read_at"`
	LastReadAt  *time.Time `db:"last_read_at"`
	ReadCount   int64      `db:"read_count"`
}

func InsertExpectedReader(dbmap *gorp.DbMap, aid, rid int64) (int64, error) {
	id, err := dbmap.SelectNullInt(`
	        with s as (
	            select id from read_receipts where article_id = $1 and reader_id = $2
	        ), i as (
	            insert into read_receipts ("article_id", "reader_id", "created_at")
	            select $1, $2, $3
	            where not exists (select 1 from s)
	            returning id
	        )
	        select id from i union all select id from s;`,
		aid,
		rid,
		time.Now().UTC(),
	)
	if err != nil {
		return -1, err
	}

	iid, err := id.Value()

	return iid.(int64), err
}

func TrackReadReceipt(dbmap *gorp.DbMap, account *Account, key, reader string) error {
	aid, err := FindArticleIdByKey(dbmap, account.Id, key)
	if err != nil {
		return err
	}

	if aid == 0 {
		aid, err = UpsertArticle(dbmap, account.Id, key)
		if err != nil {
			return err
		}
	}

	rid, err := InsertReader(dbmap, account.Id, reader)
	if err != nil {
		return err
	}

	_, err = UpsertReadReceipt(dbmap, aid, rid)
	if err != nil {
		return err
	}

	grohl.Log(grohl.Data{
		"account": account.Id,
		"reader":  reader,
		"article": key,
		"track":   "read",
	})

	return nil
}

func UpsertReadReceipt(dbmap *gorp.DbMap, articleId, readerId int64) (int64, error) {
	at := time.Now().UTC()
	_, err := dbmap.SelectNullInt(`
				update read_receipts set first_read_at = $1
				where article_id=$2 and reader_id=$3 and first_read_at is null`,
		at,
		articleId,
		readerId,
	)

	id, err := dbmap.SelectNullInt(`
        update read_receipts set last_read_at = $1
				where article_id=$2 and reader_id=$3
				returning id;
    `, at,
		articleId,
		readerId,
	)
	if err != nil {
		return -1, err
	}
	if id.Valid {
		return id.Int64, nil
	}

	id, err = dbmap.SelectNullInt(`
        with s as (
            select id from read_receipts where article_id = $1 and reader_id = $2
        ), i as (
            insert into read_receipts (
							"article_id", "reader_id",
							"created_at", "first_read_at", "last_read_at",
							"read_count")
            select $1, $2, $3, $3, $3, 0
            where not exists (select 1 from s)
            returning id
        )
        select id from i union all select id from s;
    `, articleId,
		readerId,
		at,
	)
	if err != nil {
		return 0, err
	}

	iid, err := id.Value()
	if err != nil {
		return 0, err
	}

	dbmap.Db.Exec(`update read_receipts set read_count = read_count + 1 where id = $1`, iid.(int64))

	return iid.(int64), err
}

func UnreadArticles(dbmap *gorp.DbMap, readerId int64) (keys []Article, err error) {
	_, err = dbmap.Select(&keys, `
		select articles.*, first_read_at, last_read_at
		from articles
		  inner join read_receipts on read_receipts.article_id = articles.id
		where (first_read_at is null or articles.updated_at < read_receipts.last_read_at)
		  and read_receipts.reader_id = $1`, readerId)

	return
}
