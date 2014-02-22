package readraptor

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

var dbmap *gorp.DbMap

func InitDb(connection string)  {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Account{}, "accounts").SetKeys(true, "Id")
	dbmap.AddTableWithName(ContentItem{}, "content_items").SetKeys(true, "Id")
	dbmap.AddTableWithName(Reader{}, "readers").SetKeys(true, "Id")
	dbmap.AddTableWithName(ReadReceipt{}, "read_receipts").SetKeys(true, "Id")
}
