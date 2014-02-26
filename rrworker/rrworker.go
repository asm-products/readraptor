package main

import (
	"os"
	"path/filepath"
	"time"

	rr "github.com/asm-products/readraptor/lib"
	"github.com/cupcake/gokiq"
	"github.com/garyburd/redigo/redis"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	os.Setenv("RR_ROOT", dir)

	// database
	rr.InitDb(os.Getenv("DATABASE_URL"))

	gokiq.Workers.RedisPool = redis.NewPool(rr.RedisConnect(os.Getenv("REDIS_URL")), 1)

	gokiq.Workers.PollInterval = 1 * time.Second
	gokiq.Workers.RedisNamespace = "rr"
	gokiq.Workers.WorkerCount = 5

	gokiq.Workers.Register(&rr.UserCallbackJob{})
	gokiq.Workers.Register(&rr.NewAccountEmailJob{})

	gokiq.Workers.Run()
}
