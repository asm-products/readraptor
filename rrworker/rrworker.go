package main

import (
	"os"
	"time"

	rr "github.com/asm-products/readraptor/lib"
	"github.com/whatupdave/gokiq"
)

func main() {
	// database
	rr.InitDb(os.Getenv("DATABASE_URL"))

	gokiq.Workers.PollInterval = 1 * time.Second
	gokiq.Workers.RedisNamespace = "rr"
	gokiq.Workers.WorkerCount = 200

	gokiq.Workers.Register(&rr.UserCallbackJob{})

	gokiq.Workers.Run()
}
