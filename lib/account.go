package readraptor

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/coopernurse/gorp"
)

type Account struct {
	Id       int64     `db:"id"         json:"id"`
	Created  time.Time `db:"created_at" json:"created"`
	Username string    `db:"username"   json:"username"`
	ApiKey   string    `db:"api_key"    json:"apiKey"`
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
	hasher.Write([]byte(key + username + fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
