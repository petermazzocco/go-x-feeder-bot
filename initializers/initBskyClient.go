package initializers

import (
	"context"
	"fmt"
	"os"
	"time"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

// Private variables
var (
	localClient *xrpc.Client
)

// Initialize the client once
func InitializeClient() error {
	// Only allow the initialization to happen once
	initOnce.Do(func() {
		// Context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Load envs
		host := os.Getenv("HOST")
		handle := os.Getenv("HANDLE")
		password := os.Getenv("PASSWORD")

		// Create the client
		client := &xrpc.Client{
			Host: host,
		}

		// Authenticate
		auth, err := comatproto.ServerCreateSession(ctx, client, &comatproto.ServerCreateSession_Input{
			Identifier: handle,
			Password:   password,
		})
		if err != nil {
			initError = fmt.Errorf("Error logging into Bluesky: %w", err)
			return
		}

		// Set authentication tokens for XPRC from auth response
		client.Auth = &xrpc.AuthInfo{
			AccessJwt:  auth.AccessJwt,
			RefreshJwt: auth.RefreshJwt,
			Handle:     auth.Handle,
			Did:        auth.Did,
		}

		localClient = client
	})
	return initError
}

// Get the local client
func GetClient() (*xrpc.Client, error) {
	if err := InitializeClient(); err != nil {
		return nil, err
	}

	// Allow multiple goroutines to read client simultaneously
	mutex.RLock()
	defer mutex.RUnlock()
	return localClient, nil
}

// Refresh the auth token if needed
func RefreshSession() error {
	// Only one goroutine can modify the client at a time
	mutex.Lock()
	defer mutex.Unlock()

	if localClient == nil {
		return fmt.Errorf("client not initialized")
	}
	// Create context with timeout
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return nil

}

func Shutdown() {
	mutex.Lock()
	defer mutex.Unlock()

	if localClient != nil {
		localClient = nil
	}
}
