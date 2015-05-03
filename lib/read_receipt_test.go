package readraptor_test

import (
	"testing"

	rr "github.com/asm-products/readraptor/lib"
)

func Test_TrackReadReceipt_FirstTime(t *testing.T) {
	dbmap := initTestDb(t)
	defer dbmap.Db.Close()
	a := MustCreateAccount(dbmap, "weasley@example.com")
	rid := MustCreateReader(dbmap, a.Id, "user_1")

	err := rr.TrackReadReceipt(dbmap, a, "article_1", "user_1")
	ok(t, err)

	aid, err := rr.FindArticleIdByKey(dbmap, a.Id, "article_1")
	assert(t, aid != 0, "Article not created")
	ok(t, err)

	var rec rr.ReadReceipt
	err = dbmap.SelectOne(&rec, `
		select first_read_at, last_read_at, read_count
		from read_receipts where article_id = $1 and reader_id = $2`, aid, rid)
	if err != nil {
		t.Fatal(err)
	}

	assert(t, rec.FirstReadAt != nil, "FirstReadAt should be set")
	assert(t, rec.LastReadAt != nil, "LastReadAt should be set")
	equals(t, int64(1), rec.ReadCount)
}
