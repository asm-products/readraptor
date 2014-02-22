package main

import (
	"os"
	"time"

	rr "github.com/asm-products/readraptor/lib"
	"github.com/cupcake/gokiq"
	"github.com/garyburd/redigo/redis"
)

func main() {
	// database
	rr.InitDb(os.Getenv("DATABASE_URL"))

	gokiq.Workers.RedisPool = redis.NewPool(rr.RedisConnect(os.Getenv("REDIS_URL")), 1)

	gokiq.Workers.PollInterval = 1 * time.Second
	gokiq.Workers.RedisNamespace = "rr"
	gokiq.Workers.WorkerCount = 5

	gokiq.Workers.Register(&rr.UserCallbackJob{})

	gokiq.Workers.Run()
}
