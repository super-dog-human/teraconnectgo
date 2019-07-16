package domain

import (
	"context"
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// GetAvailableGraphics for fetch graphic object from Cloud Datastore
func GetAvailableGraphics(request *http.Request) ([]Graphic, error) {
	ctx := appengine.NewContext(request)

	var graphics []Graphic

	currentUser, err := GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	usersGraphics, err := getCurrentUsersGraphics(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	graphics = append(graphics, usersGraphics...)

	publicGraphics, err := getPublicGraphics(ctx)
	if err != nil {
		return nil, err
	}
	graphics = append(graphics, publicGraphics...)

	return graphics, nil
}

func getCurrentUsersGraphics(ctx context.Context, userID string) ([]Graphic, error){
	var graphics []Graphic

	query := datastore.NewQuery("Graphic").Filter("UserId =", userID)
	keys, err := query.GetAll(ctx, &graphics)
	if err != nil {
		return nil, err
	}

	storeGraphicThumbnailUrl(ctx, &graphics, keys)

	return graphics, nil
}

func getPublicGraphics(ctx context.Context) ([]Graphic, error){
	var graphics []Graphic

	query := datastore.NewQuery("Graphic").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &graphics)
	if err != nil {
		return nil, err
	}

	err = storeGraphicThumbnailUrl(ctx, &graphics, keys)

	if err != nil {
		return nil, err
	}

	return graphics, nil
}

func storeGraphicThumbnailUrl(ctx context.Context, graphics *[]Graphic, keys []*datastore.Key) error {
	for i, key := range keys {
		id := key.StringID()
		filePath := "graphic/" + id + "." + (*graphics)[i].FileType
		fileType := "" // this is unnecessary when GET request
		bucketName := infrastructure.MaterialBucketName(ctx)
		url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)

		if err != nil {
			return err
		}

		(*graphics)[i].ID = id
		(*graphics)[i].URL = url
		(*graphics)[i].ThumbnailURL = infrastructure.GraphicThumbnailURL(ctx, id, fileType)
	}

	return nil
}