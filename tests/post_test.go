package tests

import (
	"context"
	"fmt"
	"go-x-feeder-bot/bluesky"
	"testing"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	twitterscraper "github.com/imperatrona/twitter-scraper"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testingClient *xrpc.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("An error occured loading environment variables, skipping:", err)
	}

	ctx := context.Background()

	// Create the xrpc client
	client := &xrpc.Client{
		Host: "https://bsky.social",
	}
	// Authenticate
	auth, err := comatproto.ServerCreateSession(ctx, client, &comatproto.ServerCreateSession_Input{
		Identifier: "hellopeter.dev",
		Password:   "pdxNL2020!@#",
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

	testingClient = client
}

func TestPostTextToBluesky(t *testing.T) {
	//Text
	text := "This is a test again"

	// Test posting to Bluesky with text
	if err := bluesky.PostToBluesky(text, nil, twitterscraper.Video{}, testingClient); err != nil {
		assert.Error(t, err)
		t.Fatal(err)
	}
	fmt.Println("Text posted successfully")
}

func TestPostTextToBlueskyWithImage(t *testing.T) {
	// Text
	text := "This is a test of a meme"

	// Make images
	var images []twitterscraper.Photo
	newImage := &twitterscraper.Photo{
		URL: "https://i.pinimg.com/originals/79/c1/bd/79c1bd0ea830bce7fe9f4c9e91ffa982.jpg",
		ID:  "2",
	}
	images = append(images, *newImage)
	// Test posting to Bluesky with text and images
	if err := bluesky.PostToBluesky(text, images, twitterscraper.Video{}, testingClient); err != nil {
		assert.Error(t, err)
		t.Fatal(err)
	}

	fmt.Println("Text and media posted successfully")
	assert.True(t, true)
}

func TestPostTextToBlueskyWithVideo(t *testing.T) {
	// Text
	text := "This is a test of a video"
	// Make a video
	video := &twitterscraper.Video{
		URL:     "https://video.twimg.com/amplify_video/1909279126479118337/vid/avc1/1080x1350/A18cAZZ3KqzM_O_M.mp4?tag=16",
		ID:      "1",
		Preview: "https://i.pinimg.com/originals/79/c1/bd/79c1bd0ea830bce7fe9f4c9e91ffa982.jpg",
		HLSURL:  "https://video.twimg.com/amplify_video/1909279126479118337/vid/avc1/1080x1350/A18cAZZ3KqzM_O_M.mp4?tag=16",
	}

	// Test posting to Bluesky with text and images
	if err := bluesky.PostToBluesky(text, nil, *video, testingClient); err != nil {
		assert.Error(t, err)
		t.Fatal(err)
	}

	fmt.Println("Text and media posted successfully")
	assert.True(t, true)
}
