package readraptor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cupcake/gokiq"
	"github.com/garyburd/redigo/redis"
	"github.com/technoweenie/grohl"
)

type UserCallbackJob struct {
	ReaderId int64
	Url      string
}

type UserCallback struct {
	User    string   `json:"user"`
	Pending []string `json:"pending"`
}

type UserCallbackJobEntry struct {
	Class string          `json:"class"`
	Args  UserCallbackJob `json:"args"`
}

func (j *UserCallbackJob) Perform() error {
	keys, err := UnreadArticlesMarkRead(dbmap, j.ReaderId)

	distinctId, err := dbmap.SelectStr("select distinct_id from readers where id = $1;", j.ReaderId)
	if err != nil {
		return err
	}

	callback := UserCallback{
		User:    distinctId,
		Pending: keys,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(&callback)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(j.Url, "application/json", &buf)
	if err != nil {
		panic(err)
	}

	grohl.Log(grohl.Data{
		"callback": j.Url,
		"reader":   j.ReaderId,
		"expected": keys,
		"status":   resp.Status,
	})

	return nil
}

func ScheduleCallbacks(client *gokiq.ClientConfig, readerIds []int64, at time.Time, url string) error {
	conn := client.RedisPool.Get()
	defer conn.Close()

	config := gokiq.JobConfig{
		At: at,
	}

	// Don't queue up job within 6 seconds of an existing job
	minScore := at.Add(-3 * time.Second).Unix()
	maxScore := at.Add(3 * time.Second).Unix()

	jobs, err := redis.Strings(conn.Do("ZRANGEBYSCORE", client.RedisNamespace+":schedule", minScore, maxScore))
	userJobs := make(map[int64]int64)
	for _, job := range jobs {
		var entry UserCallbackJobEntry

		if err = json.Unmarshal([]byte(job), &entry); err != nil {
			return err
		}

		userJobs[entry.Args.ReaderId] += 1
	}

	for _, rid := range readerIds {
		if userJobs[rid] > 0 {
			continue
		}

		err = client.QueueJobConfig(&UserCallbackJob{
			Url:      url,
			ReaderId: rid,
		}, config)

		if err != nil {
			return err
		}

		grohl.Log(grohl.Data{
			"schedule_callback": at,
			"url":               url,
			"reader":            rid,
		})
	}

	return nil
}

// Cribbed from https://github.com/cupcake/gokiq/blob/master/worker.go
func timeScore(t time.Time) float64 {
	return float64(t.UnixNano()) / float64(time.Second)
}
