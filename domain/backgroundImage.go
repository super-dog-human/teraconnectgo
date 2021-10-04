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
	Name   string `json:"name" datastore:",noindex"`
	URL    string `json:"url" datastore:"-"`
	SortID int64  `json:"-"`
}

type BackgroundImageErrorCode uint

const (
	BackgroundImageNotFound BackgroundImageErrorCode = 1
)

func (e BackgroundImageErrorCode) Error() string {
	switch e {
	case BackgroundImageNotFound:
		return "background image not found"
	default:
		return "unknown background image error"
	}
}

// GetAllBackgroundImagesはSortIDの照準でソートした全てのBackgroundImageを返します
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

// GetBackgroundImageは、ソートを行わずに一つだけBackGroundImageを返します
func GetBackgroundImage(ctx context.Context) (BackgroundImage, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())

	var backgroundImage BackgroundImage
	if err != nil {
		return backgroundImage, err
	}

	var images []BackgroundImage
	query := datastore.NewQuery("BackgroundImage").Limit(1)
	keys, err := client.GetAll(ctx, query, &images)
	if err != nil {
		return backgroundImage, err
	}

	if len(keys) == 0 {
		return backgroundImage, BackgroundImageNotFound
	}

	backgroundImage = images[0]
	backgroundImage.ID = keys[0].ID

	return backgroundImage, nil
}
