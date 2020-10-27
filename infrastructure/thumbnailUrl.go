package infrastructure

import (
	"context"
	"strings"
)

// AvatarThumbnailURL is return thumbnail URL of avatar each environments.
func AvatarThumbnailURL(ctx context.Context, id string) string {
	avatarThumbnailURL := storageURL(ctx) + "/avatar/{id}.png"
	return strings.Replace(avatarThumbnailURL, "{id}", id, 1)
}

// GraphicThumbnailURL is return thumbnail URL of graphic each environments.
func GraphicThumbnailURL(ctx context.Context, id string, fileType string) string {
	graphicThumbnailURL := storageURL(ctx) + "/graphic/{id}.{fileType}"
	graphicThumbnailURL = strings.Replace(graphicThumbnailURL, "{id}", id, 1)
	return strings.Replace(graphicThumbnailURL, "{fileType}", fileType, 1)

}

func storageURL(ctx context.Context) string {
	switch AppEnv() {
	case "production":
		return "https://storage.googleapis.com/teraconn_thumbnail"
	case "staging":
		return "https://storage.googleapis.com/teraconn_thumbnail_staging"
	default:
		return ""
	}
}
