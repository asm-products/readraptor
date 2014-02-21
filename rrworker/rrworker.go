package main

import (
	workers "github.com/jrallison/go-workers"
)

func main() {
	workers.Configure(map[string]string{
		"server":  "localhost:6379",
		"pool":    "30",
		"process": "1",
	})

	workers.Run()
}
