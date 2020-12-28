package domain

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// BackgroundImage type is used in the class.
type BackgroundImage struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	SortID int64  `json:"-"`
}

// GetAllBackgroundImages is return all sorted images.
func GetAllBackgroundImages(ctx context.Context) ([]BackgroundImage, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var subjects []BackgroundImage
	query := datastore.NewQuery("BackgroundImage").Order("SortID")
	keys, err := client.GetAll(ctx, query, &subjects)
	if err != nil {
		return nil, err
	}

	bucketName := infrastructure.PublicBucketName()
	for i, key := range keys {
		subjects[i].ID = key.ID
		subjects[i].URL = infrastructure.GetPublicBackGroundImageURL(bucketName, strconv.FormatInt(key.ID, 10))
	}

	return subjects, nil
}
