package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-x-feeder-bot/bluesky"
	"net/http"
	"time"

	"os"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var (
	Client   *xrpc.Client
	handle   string
	password string
)

func refreshToken(ctx context.Context, client *xrpc.Client) error {
	refresh, err := comatproto.ServerRefreshSession(ctx, client)
	if err != nil {
		auth, err := comatproto.ServerCreateSession(ctx, client, &comatproto.ServerCreateSession_Input{
			Identifier: handle,
			Password:   password,
		})
		if err != nil {
			return fmt.Errorf("failed to create new session: %w", err)
		}
		client.Auth = &xrpc.AuthInfo{
			AccessJwt:  auth.AccessJwt,
			RefreshJwt: auth.RefreshJwt,
			Handle:     auth.Handle,
			Did:        auth.Did,
		}
		return nil
	}

	client.Auth = &xrpc.AuthInfo{
		AccessJwt:  refresh.AccessJwt,
		RefreshJwt: refresh.RefreshJwt,
		Handle:     refresh.Handle,
		Did:        refresh.Did,
	}
	fmt.Println("Token refreshed successfully at:", time.Now())
	return nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("An error occured loading environment variables, skipping:", err)
	}
	// Load envs
	host := os.Getenv("HOST")
	if host == "" {
		host = "https://bsky.social" // Default value
	}

	handle = os.Getenv("HANDLE")
	if handle == "" {
		fmt.Println("Error: HANDLE environment variable is required")
		return
	}

	password = os.Getenv("PASSWORD")
	if password == "" {
		fmt.Println("Error: PASSWORD environment variable is required")
		return
	}

	// Create the xrpc client
	client := &xrpc.Client{
		Host: host,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
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
		fmt.Println("Error creating initial session:", err)
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
	if jobSpec == "" {
		jobSpec = "0 */10 * * * *" // Default value of 10 minutes
	}

	// Create a token refresh job that runs every 110 minutes (just under 2 hours)
	refreshJobSpec := "0 */110 * * * *"

	job, err := bluesky.Job(Client)
	if err != nil {
		fmt.Println("An error occured running job:", err)
		return
	}

	c := cron.New()

	// Add token refresh job
	c.AddFunc(refreshJobSpec, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := refreshToken(ctx, Client); err != nil {
			fmt.Println("Failed to refresh token:", err)
		}
	})

	// Add the main job wrapped with a token refresh check
	c.AddFunc(jobSpec, func() {
		// Check if we need to refresh token before running the job
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Run the actual job
		job()
	})

	c.Start()
	fmt.Println("Cron scheduler started. Press Ctrl+C to exit...")
	fmt.Println("XRPC Tokens will be refreshed every 110 minutes")
	select {} // This will block forever
}
