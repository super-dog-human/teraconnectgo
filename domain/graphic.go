package domain

import (
	"context"
	"net/http"
	"strings"

	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// TODO move to infrastructure for at development settings.
const graphicThumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/graphic/{id}.{fileType}"

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

		replacedURL := strings.Replace(graphicThumbnailURL, "{id}", id, 1)
		(*graphics)[i].ThumbnailURL = strings.Replace(replacedURL, "{fileType}", (*graphics)[i].FileType, 1)
	}

	return nil
}