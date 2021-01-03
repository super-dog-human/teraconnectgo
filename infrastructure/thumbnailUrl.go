package infrastructure

import (
	"context"
	"fmt"
)

// GraphicThumbnailURL is return thumbnail URL of graphic each environments.
func GraphicThumbnailURL(ctx context.Context, id int64, fileType string) string {
	return fmt.Sprintf("%s/graphic/%d.%s", storageURL(ctx), id, fileType)
}

func storageURL(ctx context.Context) string {
	return "https://storage.googleapis.com/" + ThumbnailBucketName()
}
