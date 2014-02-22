package main

import (
	"time"

	"github.com/whatupdave/gokiq"
	"github.com/asm-products/readraptor/lib"
)

func main() {
	gokiq.Workers.PollInterval = 1 * time.Second
	gokiq.Workers.RedisNamespace = "rr"
	gokiq.Workers.WorkerCount = 200

	gokiq.Workers.Register(&readraptor.UserCallbackJob{})

	gokiq.Workers.Run()
}
