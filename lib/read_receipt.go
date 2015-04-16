package readraptor

import (
	"time"

	"github.com/coopernurse/gorp"
	"github.com/technoweenie/grohl"
)

type ReadReceipt struct {
	Id         int64     `db:"id"`
	Created    time.Time `db:"created_at"`
	LastReadAt time.Time `db:"last_read_at"`
	ArticleId  int64     `db:"article_id"`
	ReaderId   int64     `db:"reader_id"`
}

func TrackReadReceipt(dbmap *gorp.DbMap, account *Account, key, reader string) error {
	cid, err := InsertArticle(dbmap, account.Id, key)
	if err != nil {
		return err
	}

	vid, err := InsertReader(dbmap, account.Id, reader)
	if err != nil {
		return err
	}

	_, err = UpsertReadReceipt(dbmap, cid, vid)
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
	id, err := dbmap.SelectNullInt(`
        update read_receipts set last_read_at = $1 where article_id=$2 and reader_id=$3 returning id;
    `, time.Now().UTC(),
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
            insert into read_receipts ("article_id", "reader_id", "created_at", "last_read_at")
            select $1, $2, $3, $3
            where not exists (select 1 from s)
            returning id
        )
        select id from i union all select id from s;
    `, articleId,
		readerId,
		time.Now().UTC(),
	)
	if err != nil {
		return -1, err
	}

	iid, err := id.Value()

	return iid.(int64), err
}

func UnreadArticles(dbmap *gorp.DbMap, readerId int64) (keys []Article, err error) {
	_, err = dbmap.Select(&keys, `
        select distinct articles.*, read_receipts.created_at as first_read_at, last_read_at from
            (select article_id from expected_readers where reader_id = $1
                except all
             select article_id from read_receipts
						   inner join articles on articles.id = read_receipts.article_id
						 where reader_id = $1 and articles.updated_at < read_receipts.last_read_at) unread_articles
        inner join articles on articles.id = unread_articles.article_id
				left join read_receipts on read_receipts.article_id = unread_articles.article_id
					and read_receipts.reader_id = $1;`, readerId)

	return
}

func UnreadArticlesMarkRead(dbmap *gorp.DbMap, readerId int64) (keys []string, err error) {
	t, err := dbmap.Begin()
	if err != nil {
		return
	}

	articles, err := UnreadArticles(dbmap, readerId)
	if err != nil {
		return
	}

	keys = make([]string, 0)
	for _, a := range articles {
		keys = append(keys, a.Key)
	}

	_, err = dbmap.Select(&keys, `
        select articles.* from
            (select article_id from expected_readers where reader_id = $1
                except all
             select article_id from read_receipts where reader_id = $1) unread_articles
        inner join articles on articles.id = unread_articles.article_id;`, readerId)

	// Mark articles as read
	for _, ci := range articles {
		_, err = UpsertReadReceipt(dbmap, ci.Id, readerId)
		if err != nil {
			return
		}
	}

	err = t.Commit()

	return
}
