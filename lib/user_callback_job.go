package readraptor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/technoweenie/grohl"
	"github.com/cupcake/gokiq"
)

type UserCallbackJob struct {
	ReaderId int64
	Url      string
}

type UserCallback struct {
	User     string   `json:"user"`
	Expected []string `json:"expected"`
}

func (j *UserCallbackJob) Perform() error {
	keys, err := UnreadContentItemKeys(dbmap, j.ReaderId)
	if err != nil {
		panic(err)
	}

	distinctId, err := dbmap.SelectStr("select distinct_id from readers where id = $1;", j.ReaderId)
	if err != nil {
		panic(err)
	}

	callback := UserCallback{
		User:     distinctId,
		Expected: keys,
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

func ScheduleCallbacks(client *gokiq.ClientConfig, readerIds []int64, callbacks []CallbackParams) {
	for _, callback := range callbacks {
		for _, rid := range readerIds {
			at := time.Now().Add(time.Duration(callback.Seconds) * time.Second)

			config := gokiq.JobConfig{
				At: at,
			}

			client.QueueJobConfig(&UserCallbackJob{
				Url:      callback.Url,
				ReaderId: rid,
			}, config)

			grohl.Log(grohl.Data{
				"schedule_job": at,
				"url":          callback.Url,
				"reader":       rid,
			})
		}
	}
}
