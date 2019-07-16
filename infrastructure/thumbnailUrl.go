package infrastructure

import (
	"context"
	"strings"

	"google.golang.org/appengine"
)

// AvatarThumbnailURL is return thumbnail URL of avatar each environments.
func AvatarThumbnailURL(ctx context.Context, id string) string {
	avatarThumbnailURL := storageUrl(ctx) + "/avatar/{id}.png"
	return strings.Replace(avatarThumbnailURL, "{id}", id, 1)
}

// GraphicThumbnailURL is return thumbnail URL of graphic each environments.
func GraphicThumbnailURL(ctx context.Context, id string, fileType string) string {
	graphicThumbnailURL := storageUrl(ctx) + "/graphic/{id}.{fileType}"
	graphicThumbnailURL = strings.Replace(graphicThumbnailURL, "{id}", id, 1)
	return strings.Replace(graphicThumbnailURL, "{fileType}", fileType, 1)

}

func storageUrl(ctx context.Context) string {
	storageUrl := "https://storage.googleapis.com/teraconn_thumbnail"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return storageUrl
	}

	return storageUrl + "_development"
}