package bluesky

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"strings"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func removeURLs(text string) string {
	// Twitter scrapper returns a https://pbs.twimg/ in the Text
	// We can remove it here and just post the text
	urlPattern := regexp.MustCompile(`https?://\S+`)
	return strings.TrimSpace(urlPattern.ReplaceAllString(text, ""))
}

// PostToBluesky posts text and media to Bluesky
func PostToBluesky(text string, images []twitterscraper.Photo, video twitterscraper.Video, client *xrpc.Client) error {
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Remove URL from text
	cleanedText := removeURLs(text)
	fmt.Println("Incoming post content via Twitter:", cleanedText)

	// Grab the feed from the author to check if the cleanedText already exists
	feed, err := bsky.FeedGetAuthorFeed(ctx, client, client.Auth.Did, "", "", false, 10)
	if err != nil {
		fmt.Println("Error fetching feed:", err)
		return err
	}

	// Loop through the feed and check if the cleanedText already exists
	for _, post := range feed.Feed {
		postRecord := post.Post.Record.Val.(*bsky.FeedPost)
		fmt.Printf("Post: %s\n", postRecord.Text)
		if postRecord.Text == cleanedText {
			fmt.Println("Post already exists, skipping")
			return nil
		}
	}

	// Create image blobs if available
	var imageBlobs []*bsky.EmbedImages_Image
	if len(images) > 0 {
		for _, image := range images {
			// Check if image URL is valid
			if image.URL == "" {
				fmt.Println("Empty image URL, skipping")
				continue
			}

			imageBlob, err := UploadBlobToRepo(ctx, client, image.URL)
			if err != nil {
				fmt.Println("An error occurred while uploading image:", err)
				continue
			}

			embeddedImage := &bsky.EmbedImages_Image{
				Alt:   "This X feeder account posted an image",
				Image: imageBlob,
			}
			imageBlobs = append(imageBlobs, embeddedImage)
		}
	}

	// Create video blob if available
	var videoBlob *util.LexBlob
	if video.URL != "" {
		blob, err := UploadBlobToRepo(ctx, client, video.URL)
		if err != nil {
			fmt.Println("An error occurred while uploading video:", err)
			// Continue execution, don't return
		} else {
			videoBlob = blob
		}
	}

	// Create the post record
	record := &bsky.FeedPost{
		LexiconTypeID: "app.bsky.feed.post",
		Text:          cleanedText,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		Langs:         []string{"en"},
	}

	// Add embeds conditionally
	if len(imageBlobs) > 0 || videoBlob != nil {
		record.Embed = &bsky.FeedPost_Embed{}

		if len(imageBlobs) > 0 {
			record.Embed.EmbedImages = &bsky.EmbedImages{
				LexiconTypeID: "app.bsky.embed.images",
				Images:        imageBlobs,
			}
		}

		if videoBlob != nil {
			record.Embed.EmbedVideo = &bsky.EmbedVideo{
				LexiconTypeID: "app.bsky.embed.video",
				Video:         videoBlob,
			}
		}
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
		return err
	}

	// Print URI for verification
	fmt.Println("Post URI:", resp.Uri)
	return nil
}
