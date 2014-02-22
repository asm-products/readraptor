package readraptor

import (
	"time"

	"github.com/technoweenie/grohl"
	"github.com/whatupdave/gokiq"
)

type UserCallbackJob struct {
	ReaderId int64
	Url      string
}

func (j *UserCallbackJob) Perform() error {
	grohl.Log(grohl.Data{
		"callback": j.Url,
		"reader":   j.ReaderId,
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
