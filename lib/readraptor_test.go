package readraptor_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	rr "github.com/asm-products/readraptor/lib"
	"github.com/coopernurse/gorp"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v (%T)\n\n\tgot: %#v (%T)\033[39m\n\n", filepath.Base(file), line, exp, exp, act, act)
		tb.FailNow()
	}
}

func expectInclude(tb testing.TB, a, b string) {
	if strings.Index(a, b) == -1 {
		tb.Errorf("Expected '%v' to include '%v'", a, b)
	}
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func MustCreateAccount(dbmap *gorp.DbMap, email string) *rr.Account {
	account := rr.NewAccount(email)
	err := dbmap.Insert(account)
	if err != nil {
		panic(err.Error())
	}
	return account
}

func MustCreateReader(dbmap *gorp.DbMap, accountId int64, distinctId string) int64 {
	rid, err := rr.InsertReader(dbmap, accountId, distinctId)
	if err != nil {
		panic(err.Error())
	}
	return rid
}

func initTestDb(t *testing.T) *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	dbmap := rr.InitDb("postgres://localhost/rr_test?sslmode=disable")

	err := dbmap.CreateTablesIfNotExists()
	ok(t, err)

	// delete any existing rows
	err = dbmap.TruncateTables()
	ok(t, err)

	return dbmap
}
