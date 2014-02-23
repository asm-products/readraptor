package main

import (
	"os"
	"path/filepath"

	rr "github.com/asm-products/readraptor/lib"
	"github.com/coopernurse/gorp"
)

var dbmap *gorp.DbMap

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	rr.RunWeb(dir)
}
