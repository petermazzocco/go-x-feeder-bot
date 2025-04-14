package bluesky

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

func UploadBlobToRepo(ctx context.Context, client *xrpc.Client, url string) (*util.LexBlob, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is empty meaning no items to upload, skipping...")
	}

	// Download item
	fmt.Println("Downloading item...")
	item, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to get item", err)
		return nil, err
	}
	defer item.Body.Close()

	// Upload the image directly from the response body to the repo via atproto (not bsky)
	fmt.Println("Uploading item...")
	uploadResp, err := atproto.RepoUploadBlob(ctx, client, item.Body)
	if err != nil {
		fmt.Println("Failed to upload blob", err)
		return nil, err
	}

	// Create the LexBlob with the reference from the upload
	blob := &util.LexBlob{
		Ref:      uploadResp.Blob.Ref,
		MimeType: item.Header.Get("Content-Type"),
		Size:     uploadResp.Blob.Size,
	}
	fmt.Println("Item uploaded successfully", blob)
	return blob, nil
}
