package main

import (
	"fmt"
	"go-x-feeder-bot/cronJobs"
	"go-x-feeder-bot/initializers"

	"os"

	"github.com/robfig/cron"
)

func init() {
	if err := initializers.InitializeEVN(); err != nil {
		fmt.Println("An error occured initializing environment variables:", err)
	}

	if err := initializers.InitializeClient(); err != nil {
		fmt.Println("An error occured initializing client:", err)
	}

	if err := initializers.InitializeScrapper(); err != nil {
		fmt.Println("An error occured initializing scrapper:", err)
	}
}

func main() {
	jobSpec := os.Getenv("JOB_SPEC")

	job, err := cronJobs.Job()
	if err != nil {
		fmt.Println("An error occured running job:", err)
		return
	}

	// Instantiate new cron job
	c := cron.New()

	c.AddFunc(jobSpec, job)
	c.Start()
}
