package bluesky

import (
	"context"
	"fmt"
	"os"
	"testing"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	twitterscraper "github.com/imperatrona/twitter-scraper"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var client *xrpc.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("An error occured loading environment variables, skipping:", err)
	}
	handle := os.Getenv("HANDLE")
	if handle == "" {
		fmt.Println("Error: HANDLE environment variable is required")
		return
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		fmt.Println("Error: PASSWORD environment variable is required")
		return
	}
	ctx := context.Background()
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
}

func TestPostTextToBluesky(t *testing.T) {
	// Skip test in CI environment or when running automated tests
	t.Skip("Skipping test that requires authentication")

	//Text
	text := "This is a test"

	// Test posting to Bluesky with text
	if err := PostToBluesky(text, nil, twitterscraper.Video{}, client); err != nil {
		assert.Error(t, err)
	}
	fmt.Println("Text posted successfully")
}

func TestPostTextToBlueskyWithMedia(t *testing.T) {
	// Skip test in CI environment or when running automated tests
	t.Skip("Skipping test that requires authentication")

	// Text
	text := "This is a test"

	// Make images
	var images []twitterscraper.Photo
	newImage := &twitterscraper.Photo{
		URL: "http://example.com/image.jpg",
		ID:  "2",
	}
	images = append(images, *newImage)

	// Make a video
	video := &twitterscraper.Video{
		URL:     "http://example.com/video.mp4",
		ID:      "1",
		Preview: "http://example.com/image.jpg",
		HLSURL:  "http://example.com/video.mp4",
	}

	// Test posting to Bluesky with text and images
	if err := PostToBluesky(text, images, *video, client); err != nil {
		assert.Error(t, err)
	}

	fmt.Println("Text and media posted successfully")
}
