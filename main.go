package main

import (
	"fmt"
	"go-x-feeder-bot/bluesky"

	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("An error occured loading environment variables, skipping:", err)
	}
}

func main() {
	jobSpec := os.Getenv("JOB_SPEC")

	job, err := bluesky.Job()
	if err != nil {
		fmt.Println("An error occured running job:", err)
		return
	}

	// Instantiate new cron job
	c := cron.New()

	job()

	c.AddFunc(jobSpec, job)
	c.Start()
}
