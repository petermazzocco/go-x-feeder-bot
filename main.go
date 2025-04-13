package main

import (
	"context"
	"fmt"
	"go-x-feeder-bot/bluesky"
	"time"

	"os"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var Client *xrpc.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("An error occured loading environment variables, skipping:", err)
	}
	// Load envs
	host := os.Getenv("HOST")
	handle := os.Getenv("HANDLE")
	password := os.Getenv("PASSWORD")

	// Create the xrpc client
	client := &xrpc.Client{
		Host: host,
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Authenticate
	auth, err := comatproto.ServerCreateSession(ctx, client, &comatproto.ServerCreateSession_Input{
		Identifier: handle,
		Password:   password,
	})
	if err != nil {
		fmt.Println("An error occurred logging into Bluesky:", err)
		return
	}

	// Set authentication tokens for XPRC from auth response
	client.Auth = &xrpc.AuthInfo{
		AccessJwt:  auth.AccessJwt,
		RefreshJwt: auth.RefreshJwt,
		Handle:     auth.Handle,
		Did:        auth.Did,
	}

	// set client var
	Client = client
}

func main() {
	jobSpec := os.Getenv("JOB_SPEC")

	job, err := bluesky.Job(Client)
	if err != nil {
		fmt.Println("An error occured running job:", err)
		return
	}

	c := cron.New()
	c.AddFunc(jobSpec, job)
	c.Start()
}
