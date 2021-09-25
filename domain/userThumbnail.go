package domain

import (
	"context"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// CreateUserThumbnailBlankFile is create blank image file to public bucket.
func CreateUserThumbnailBlankFile(ctx context.Context, id int64) (string, error) {
	fileDir := "user"
	fileName := strconv.FormatInt(id, 10)
	fileRequest := infrastructure.FileRequest{
		Extension:   "png",
		ContentType: "image/png",
	}

	var url string
	var err error

	url, err = infrastructure.CreateBlankFileToPublicGCS(ctx, fileName, fileDir, fileRequest)
	if err != nil {
		return "", err
	}

	return url, nil
}
