package readraptor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cupcake/gokiq"
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

func (j *UserCallbackJob) Perform() error {
	articles, err := UnreadArticles(dbmap, j.ReaderId)
	if err != nil {
		panic(err)
	}

	distinctId, err := dbmap.SelectStr("select distinct_id from readers where id = $1;", j.ReaderId)
	if err != nil {
		panic(err)
	}

	keys := make([]string, 0)
	for _, a := range articles {
		keys = append(keys, a.Key)
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

	// Mark articles as read
	for _, ci := range articles {
		_, err := InsertReadReceipt(dbmap, ci.Id, j.ReaderId)
		if err != nil {
			panic(err)
		}
	}

	grohl.Log(grohl.Data{
		"callback": j.Url,
		"reader":   j.ReaderId,
		"expected": keys,
		"status":   resp.Status,
	})

	return nil
}

func ScheduleCallbacks(client *gokiq.ClientConfig, readerIds []int64, at time.Time, url string) {
	config := gokiq.JobConfig{
		At: at,
	}

	for _, rid := range readerIds {
		client.QueueJobConfig(&UserCallbackJob{
			Url:      url,
			ReaderId: rid,
		}, config)

		grohl.Log(grohl.Data{
			"schedule_job": at,
			"url":          url,
			"reader":       rid,
		})
	}
}
