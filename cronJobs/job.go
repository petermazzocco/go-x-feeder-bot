package cronJobs

import (
	"context"
	"fmt"
	"go-x-feeder-bot/bluesky"
	"go-x-feeder-bot/initializers"
	"os"
	"time"
)

func Job() (func(), error) {
	scraper, err := initializers.GetScrapper()
	if err != nil {
		fmt.Println("An error occured initializing scrapper:", err)
		return nil, err
	}
	twitterAccount := os.Getenv("TWITTER_ACCOUNT")

	fmt.Println("Check Twitter at:", time.Now())
	for tweet := range scraper.GetTweets(context.Background(), twitterAccount, 1) {
		if tweet.Error != nil {
			fmt.Println("An error occured fetching tweets:", tweet.Error)
			return nil, fmt.Errorf("error fetching tweets: %w", tweet.Error)
		}
		bluesky.PostToBluesky(tweet.Text, tweet.Photos, tweet.Videos[0])
	}
	return nil, nil
}
