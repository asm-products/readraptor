package readraptor

import (
	"testing"
	"time"

	"github.com/cupcake/gokiq"
	"github.com/garyburd/redigo/redis"
)

func Test_ScheduleCallbacks(t *testing.T) {
	pool := redis.NewPool(RedisConnect("redis://localhost:6379/6"), 1)
	gokiq.Client.RedisPool = pool
	gokiq.Client.RedisNamespace = "test"

	conn := pool.Get()
	defer conn.Close()

	conn.Do("del", "test:schedule")

	rids := []int64{1, 2}
	err := ScheduleCallbacks(gokiq.Client, rids, time.Now().UTC(), "http://example.com/webhook")
	if err != nil {
		t.Fatal(err)
	}
	err = ScheduleCallbacks(gokiq.Client, rids, time.Now().UTC(), "http://example.com/webhook")
	if err != nil {
		t.Fatal(err)
	}

	jobs, _ := conn.Do("ZCARD", "test:schedule")

	expect(t, int64(2), jobs)
}
