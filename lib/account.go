package readraptor

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	stripe "github.com/bradrydzewski/go.stripe"
	"github.com/cupcake/gokiq"
	"github.com/martini-contrib/sessionauth"
	"github.com/technoweenie/grohl"
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
	ConfirmedAt        *time.Time `db:"confirmed_at"         json:"-"`
	// Stripe
	CustomerId *string `db:"customer_id"  json:"-"`

	// Session
	authenticated bool `db:"-" json:"-"`
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
	fmt.Println("  Tried", account)
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

// Session
func GuestAccount() sessionauth.User {
	return &Account{}
}

func (a *Account) IsAuthenticated() bool {
	return a.authenticated
}

func (a *Account) Login() {
	a.authenticated = true
}

func (a *Account) Logout() {
	a.authenticated = true
}

func (a *Account) UniqueId() interface{} {
	return a.Id
}

func (a *Account) GetById(id interface{}) error {
	return dbmap.SelectOne(a, "select * from accounts where id = $1", id)
}

// Stripe
func (a *Account) CreateStripeCustomer(token string) error {
	stripe.SetKey(os.Getenv("STRIPE_SECRET"))

	params := stripe.CustomerParams{
		Email: a.Email,
		Token: token,
	}

	customer, err := stripe.Customers.Create(&params)
	if err != nil {
		return err
	}

	_, err = dbmap.Exec(
		`update accounts set customer_id = $2 where id = $1`,
		a.Id,
		customer.Id,
	)
	return err
}

func genKey(input string) string {
	salt := os.Getenv("API_GEN_SECRET") + fmt.Sprintf("%d", time.Now().UnixNano())
	hasher := sha1.New()
	hasher.Write([]byte(salt + input))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
