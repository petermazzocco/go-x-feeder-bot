package bluesky

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	if authToken == "" {
		fmt.Println("Error: AUTH TOKEN environment variable is required")
		return nil, fmt.Errorf("AUTH TOKEN environment variable is required")

	}
	csrfToken := os.Getenv("CSRF_TOKEN")
	if csrfToken == "" {
		fmt.Println("Error: CSRF TOKEN environment variable is required")
		return nil, fmt.Errorf("CSRF TOKEN environment variable is required")
	}
	twitterAccount := os.Getenv("TWITTER_ACCOUNT")
	if twitterAccount == "" {
		fmt.Println("Error: TWITTER ACCOUNT environment variable is required")
		return nil, fmt.Errorf("TWITTER ACCOUNT environment variable is required")
	}

	// Return a function, not nil
	return func() {
		scraper := twitterscraper.New()
		scraper.SetAuthToken(twitterscraper.AuthToken{Token: authToken, CSRFToken: csrfToken})

		// After setting Cookies or AuthToken you have to execute IsLoggedIn method.
		// Add debug logging
		fmt.Printf("Auth token length: %d, starts with: %s...\n", len(authToken), authToken[:5])
		fmt.Printf("CSRF token length: %d, starts with: %s...\n", len(csrfToken), csrfToken[:5])

		if !scraper.IsLoggedIn() {
			fmt.Println("Invalid auth tokens - cannot authenticate with Twitter")
			// Additional debugging info
			resp, err := scraper.GetProfile("beaverfootball")
			if err != nil {
				fmt.Printf("Error info: %v\n", err)
			} else {
				fmt.Printf("Got response but still not logged in, status: %v\n", resp)
			}
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

			// Try to post to Bluesky with error handling for auth errors
			err := PostToBluesky(tweet.Text, photos, video, client)
			if err != nil {
				fmt.Println("Error posting to Bluesky:", err)
				
				// Check if the error might be an authentication error
				// This is a simplified check - you might want to make this more specific 
				// based on the actual error format returned by the API
				if isAuthError(err) {
					fmt.Println("Authentication error detected, token may have expired")
				}
			}
		}
	}, nil
}

// isAuthError checks if an error is likely related to authentication
func isAuthError(err error) bool {
	errStr := err.Error()
	// Check for common auth-related error strings
	authErrorIndicators := []string{
		"unauthorized", 
		"Unauthorized",
		"401",
		"auth",
		"token",
		"expired",
		"authentication",
		"login",
		"credentials",
		"jwt",
	}
	
	for _, indicator := range authErrorIndicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}
	
	return false
}
