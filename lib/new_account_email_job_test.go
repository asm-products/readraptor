package readraptor_test

import (
	"os"
	"testing"

	rr "github.com/asm-products/readraptor/lib"
	_ "github.com/lib/pq"
)

func Test_NewAccountEmailJob(t *testing.T) {
	dbmap := initTestDb(t)
	defer dbmap.Db.Close()

	os.Setenv("RR_ROOT", "..")

	account := rr.NewAccount("joe@crabshack.com")
	token := "confirm1234"
	account.ConfirmationToken = &token
	err := dbmap.Insert(account)
	ok(t, err)

	job := rr.NewAccountEmailJob{
		AccountId: account.Id,
	}

	message, err := job.CreateMessage(account)
	ok(t, err)

	expectInclude(t, message.Body, "/confirm/confirm1234")
}
