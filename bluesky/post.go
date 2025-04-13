package bluesky

import (
	"context"
	"fmt"
	"go-x-feeder-bot/initializers"
	"regexp"
	"time"

	"strings"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func removeURLs(text string) string {
	// Twitter scrapper returns a https://pbs.twimg/ in the Text
	// We can remove it here and just post the text
	urlPattern := regexp.MustCompile(`https?://\S+`)
	return strings.TrimSpace(urlPattern.ReplaceAllString(text, ""))
}

// Log into Bluesky
func PostToBluesky(text string, images []twitterscraper.Photo, video twitterscraper.Video) {
	client, err := initializers.GetClient()
	if err != nil {
		fmt.Println("An error occured getting the client", err)
	}
	ctx := context.Background()

	// Remove URL from text
	cleanedText := removeURLs(text)
	fmt.Println("Incoming post content via Twitter", cleanedText)

	// Grab the feed from the author to check if the cleanedText already exists
	feed, err := bsky.FeedGetAuthorFeed(ctx, client, client.Auth.Did, "", "", false, 10)
	if err != nil {
		fmt.Println("Error fetching feed:", err)
		return
	}

	// Loop through the feed and check if the cleanedText already exists
	for _, post := range feed.Feed {
		postRecord := post.Post.Record.Val.(*bsky.FeedPost)
		fmt.Printf("Post: %s\n", postRecord.Text)
	}

	// Create image blobs if available
	var imageBlobs []*bsky.EmbedImages_Image
	for _, image := range images {
		imageBlob, err := UploadBlobToRepo(ctx, client, image.URL)
		if err != nil {
			fmt.Println("An error occurred while uploading image")
			continue
		}
		embeddedImage := &bsky.EmbedImages_Image{
			Alt:   "This X feeder account posted an image",
			Image: imageBlob,
		}
		imageBlobs = append(imageBlobs, embeddedImage)
	}

	// Create video blob if available
	videoBlob, err := UploadBlobToRepo(ctx, client, video.URL)
	if err != nil {
		fmt.Println("An error occurred while uploading video")
		return
	}
	// Check if the image blobs or video blob is not nil
	if imageBlobs != nil || videoBlob != nil {
		// Create the post record with the blob
		record := &bsky.FeedPost{
			LexiconTypeID: "app.bsky.feed.post",
			Text:          cleanedText,
			CreatedAt:     time.Now().UTC().Format(time.RFC3339),
			Langs:         []string{"en"},
			Embed: &bsky.FeedPost_Embed{
				EmbedImages: &bsky.EmbedImages{
					LexiconTypeID: "app.bsky.embed.images",
					Images:        imageBlobs,
				},
				EmbedVideo: &bsky.EmbedVideo{
					LexiconTypeID: "app.bsky.embed.video",
					Video:         videoBlob,
				},
			},
		}

		// Submit the post to the Bluesky network and to the repo
		resp, err := comatproto.RepoCreateRecord(ctx, client, &comatproto.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Repo:       client.Auth.Did,
			Record: &lexutil.LexiconTypeDecoder{
				Val: record,
			},
		})
		if err != nil {
			fmt.Println("Error creating post:", err)
			return
		}

		// Print uri for verification
		fmt.Println("Post URI", resp.Uri)
	} else {
		// Create the post record without a blob
		record := &bsky.FeedPost{
			LexiconTypeID: "app.bsky.feed.post",
			Text:          cleanedText,
			CreatedAt:     time.Now().UTC().Format(time.RFC3339),
			Langs:         []string{"en"},
		}

		// Submit the post to the Bluesky network and to the repo
		resp, err := comatproto.RepoCreateRecord(ctx, client, &comatproto.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Repo:       client.Auth.Did,
			Record: &lexutil.LexiconTypeDecoder{
				Val: record,
			},
		})
		if err != nil {
			fmt.Println("Error creating post:", err)
			return
		}

		// Print uri for verification
		fmt.Println("Post URI", resp.Uri)
	}
}
