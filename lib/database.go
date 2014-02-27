package readraptor

import (
	"database/sql"
	"net/url"

	"github.com/coopernurse/gorp"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
)

var dbmap *gorp.DbMap

func InitDb(connection string) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Account{}, "accounts").SetKeys(true, "Id")
	dbmap.AddTableWithName(Article{}, "articles").SetKeys(true, "Id")
	dbmap.AddTableWithName(Reader{}, "readers").SetKeys(true, "Id")
	dbmap.AddTableWithName(ReadReceipt{}, "read_receipts").SetKeys(true, "Id")

	// dbmap.TraceOn("[gorp]", log.New(os.Stdout, "sql:", log.Lmicroseconds))
}

func RedisConnect(connection string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		url, err := url.Parse(connection)
		if err != nil {
			return nil, err
		}

		c, err := redis.Dial("tcp", url.Host)
		if err != nil {
			return nil, err
		}

		if url.User != nil {
			password, set := url.User.Password()

			if set {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
		}
		return c, err
	}
}
