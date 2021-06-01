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

	var images []BackgroundImage
	query := datastore.NewQuery("BackgroundImage").Order("SortID")
	keys, err := client.GetAll(ctx, query, &images)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		images[i].ID = key.ID
		images[i].URL = infrastructure.GetPublicBackgroundImageURL(strconv.FormatInt(key.ID, 10))
	}

	return images, nil
}
