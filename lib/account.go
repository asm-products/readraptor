package readraptor

import (
	"crypto/sha1"
	"fmt"
	"github.com/cupcake/gokiq"
	"github.com/technoweenie/grohl"
	"os"
	"time"
)

type Account struct {
	Id         int64     `db:"id"          json:"id"`
	Created    time.Time `db:"created_at"  json:"created"`
	Email      string    `db:"email"       json:"email"`
	PublicKey  string    `db:"public_key"  json:"-"`
	PrivateKey string    `db:"private_key" json:"-"`

	// Confirmable
	ConfirmationToken  *string    `db:"confirmation_token"   json:"-"`
	ConfirmationSentAt *time.Time `db:"confirmation_sent_at" json:"-"`
	ConfirmedAt        *string    `db:"confirmed_at"         json:"-"`
}

func NewAccount(email string) *Account {
	return &Account{
		Created:    time.Now(),
		Email:      email,
		PublicKey:  genKey("public" + email),
		PrivateKey: genKey("private" + email),
	}
}

func FindAccount(id int64) (*Account, error) {
	return FindAccountBy("id", id)
}

func FindAccountByPublicKey(key string) (*Account, error) {
	return FindAccountBy("public_key", key)
}

func FindAccountByConfirmationToken(token string) (*Account, error) {
	return FindAccountBy("confirmation_token", token)
}

func FindAccountBy(column string, value interface{}) (*Account, error) {
	var account Account
	err := dbmap.SelectOne(&account, "select * from accounts where "+column+" = $1", value)
	return &account, err
}

func (a *Account) SendNewAccountEmail(client *gokiq.ClientConfig) {
	client.QueueJob(&NewAccountEmailJob{
		AccountId: a.Id,
	})

	grohl.Log(grohl.Data{
		"queue":   "NewAccountEmailJob",
		"account": a.Id,
	})
}

func genKey(input string) string {
	salt := os.Getenv("API_GEN_SECRET") + fmt.Sprintf("%d", time.Now().UnixNano())
	hasher := sha1.New()
	hasher.Write([]byte(salt + input))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
