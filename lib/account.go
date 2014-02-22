package readraptor

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/coopernurse/gorp"
)

type Account struct {
	Id       int64     `db:"id"`
	Created  time.Time `db:"created_at"`
	Username string    `db:"username"`
	ApiKey   string    `db:"api_key"`
}

func NewAccount(username string) *Account {
	return &Account{
		Username: username,
		Created:  time.Now(),
		ApiKey:   genApiKey(username),
	}
}

func FindAccountByUsername(dbmap *gorp.DbMap, username string) (*Account, error) {
	var account Account
	err := dbmap.SelectOne(&account, "select * from accounts where username = $1", username)
	return &account, err
}

func genApiKey(username string) string {
	key := os.Getenv("API_GEN_SECRET")
	hasher := sha1.New()
	hasher.Write([]byte(key + username))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
