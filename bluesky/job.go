package bluesky

import (
	"context"
	"fmt"
	"os"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func Job() (func(), error) {
	// Get ENV variables
	authToken := os.Getenv("AUTH_TOKEN")
	csrfToken := os.Getenv("CSRF_TOKEN")

	scraper := twitterscraper.New()
	scraper.SetAuthToken(twitterscraper.AuthToken{Token: authToken, CSRFToken: csrfToken})

	// After setting Cookies or AuthToken you have to execute IsLoggedIn method.
	// Without it, scraper wouldn't be able to make requests that requires authentication
	if !scraper.IsLoggedIn() {
		panic("Invalid AuthToken")
	}

	twitterAccount := os.Getenv("TWITTER_ACCOUNT")

	fmt.Println("Check Twitter at:", time.Now())
	for tweet := range scraper.GetTweets(context.Background(), twitterAccount, 1) {
		fmt.Print("Getting tweets\n")
		if tweet.Error != nil {
			fmt.Println("An error occurred fetching tweets:", tweet.Error)
			return nil, fmt.Errorf("error fetching tweets: %w", tweet.Error)
		}

		var video twitterscraper.Video
		if len(tweet.Videos) > 0 {
			video = tweet.Videos[0]
		}

		var photos []twitterscraper.Photo
		if len(tweet.Photos) > 0 {
			photos = tweet.Photos
		}

		PostToBluesky(tweet.Text, photos, video)
	}
	return nil, nil
}
