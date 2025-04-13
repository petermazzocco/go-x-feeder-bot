package bluesky

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bluesky-social/indigo/xrpc"
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func Job(client *xrpc.Client) (func(), error) {
	// Check if client is nil
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Get ENV variables
	authToken := os.Getenv("AUTH_TOKEN")
	csrfToken := os.Getenv("CSRF_TOKEN")
	twitterAccount := os.Getenv("TWITTER_ACCOUNT")

	// Return a function, not nil
	return func() {
		scraper := twitterscraper.New()
		scraper.SetAuthToken(twitterscraper.AuthToken{Token: authToken, CSRFToken: csrfToken})

		// After setting Cookies or AuthToken you have to execute IsLoggedIn method.
		if !scraper.IsLoggedIn() {
			fmt.Println("Invalid auth tokens")
			return
		}

		fmt.Println("Check Twitter at:", time.Now())
		for tweet := range scraper.GetTweets(context.Background(), twitterAccount, 1) {
			fmt.Print("Getting tweets\n")
			if tweet.Error != nil {
				fmt.Println("An error occurred fetching tweets:", tweet.Error)
				return
			}

			// Check if there are videos or photos
			var video twitterscraper.Video
			if len(tweet.Videos) > 0 {
				video = tweet.Videos[0]
			}

			var photos []twitterscraper.Photo
			if len(tweet.Photos) > 0 {
				photos = tweet.Photos
			}

			err := PostToBluesky(tweet.Text, photos, video, client)
			if err != nil {
				fmt.Println("Error posting to Bluesky:", err)
			}
		}
	}, nil
}
