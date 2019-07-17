package domain

import (
	"context"

	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"google.golang.org/appengine/datastore"
)

func GetGraphicsByIds(ctx context.Context, ids []string) ([]Graphic, error) {
	var graphicKeys []*datastore.Key

	for _, id := range ids {
		graphicKeys = append(graphicKeys, datastore.NewKey(ctx, "Graphic", id, 0, nil))
	}

	graphics := make([]Graphic, len(ids))
	if err := datastore.GetMulti(ctx, graphicKeys, graphics); err != nil {
		return nil, err
	}

	for i, id := range ids {
		graphics[i].ID = id
	}

	return graphics, nil
}

func GetCurrentUsersGraphics(ctx context.Context, userID string) ([]Graphic, error){
	var graphics []Graphic

	query := datastore.NewQuery("Graphic").Filter("UserId =", userID)
	keys, err := query.GetAll(ctx, &graphics)
	if err != nil {
		return nil, err
	}

	storeGraphicThumbnailUrl(ctx, &graphics, keys)

	return graphics, nil
}

func GetPublicGraphics(ctx context.Context) ([]Graphic, error){
	var graphics []Graphic

	query := datastore.NewQuery("Graphic").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &graphics)
	if err != nil {
		return nil, err
	}

	if err = storeGraphicThumbnailUrl(ctx, &graphics, keys); err != nil {
		return nil, err
	}

	return graphics, nil
}

func GetGraphicFileTypes(ctx context.Context, graphicIDs []string) (map[string]string, error) {
	var keys []*datastore.Key
	for _, id := range graphicIDs {
		keys = append(keys, datastore.NewKey(ctx, "Graphic", id, 0, nil))
	}

	graphicFileTypes := map[string]string{}
	graphics := make([]Graphic, len(graphicIDs))
	if err := datastore.GetMulti(ctx, keys, graphics); err != nil {
		return nil, err
	}

	for i, g := range graphics {
		id := graphicIDs[i]
		graphicFileTypes[id] = g.FileType
	}
	return graphicFileTypes, nil
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
