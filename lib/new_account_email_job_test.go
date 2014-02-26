package readraptor

import (
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func Test_NewAccountEmailJob(t *testing.T) {
	// initialize the DbMap
	dbmap := initTestDb(t)
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(t, err, "TruncateTables failed")

	os.Setenv("RR_ROOT", "..")

	account := NewAccount("joe@crabshack.com")
	token := "confirm1234"
	account.ConfirmationToken = &token
	err = dbmap.Insert(account)
	checkErr(t, err, "Insert failed")

	job := NewAccountEmailJob{
		AccountId: account.Id,
	}

	message, err := job.CreateMessage(account)
	checkErr(t, err, "Job failed")

	expectInclude(t, message.Body, "/confirm/confirm1234")
}

func expectInclude(t *testing.T, a, b string) {
	if strings.Index(a, b) == -1 {
		t.Errorf("Expected '%v' to include '%v'", a, b)
	}
}
