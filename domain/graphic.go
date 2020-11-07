package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetGraphicsByIDs(ctx context.Context, ids []int64) ([]Graphic, error) {
	var graphicKeys []*datastore.Key

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		graphicKeys = append(graphicKeys, datastore.IDKey("Graphic", id, nil))
	}

	graphics := make([]Graphic, len(ids))
	if err := client.GetMulti(ctx, graphicKeys, graphics); err != nil {
		return nil, err
	}

	for i, id := range ids {
		graphics[i].ID = id
	}

	return graphics, nil
}

func GetCurrentUsersGraphics(ctx context.Context, userID int64) ([]Graphic, error) {
	var graphics []Graphic

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("Graphic").Filter("UserID =", userID)
	keys, err := client.GetAll(ctx, query, &graphics)
	if err != nil {
		return nil, err
	}

	storeGraphicThumbnailURL(ctx, &graphics, keys)

	return graphics, nil
}

func GetPublicGraphics(ctx context.Context) ([]Graphic, error) {
	var graphics []Graphic

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("Graphic").Filter("IsPublic =", true)
	keys, err := client.GetAll(ctx, query, &graphics)
	if err != nil {
		return nil, err
	}

	if err = storeGraphicThumbnailURL(ctx, &graphics, keys); err != nil {
		return nil, err
	}

	return graphics, nil
}

func GetGraphicFileTypes(ctx context.Context, graphicIDs []int64) (map[int64]string, error) {
	var keys []*datastore.Key

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	for _, id := range graphicIDs {
		keys = append(keys, datastore.IDKey("Graphic", id, nil))
	}

	graphicFileTypes := map[int64]string{}
	graphics := make([]Graphic, len(graphicIDs))
	if err := client.GetMulti(ctx, keys, graphics); err != nil {
		return nil, err
	}

	for i, g := range graphics {
		id := graphicIDs[i]
		graphicFileTypes[id] = g.FileType
	}
	return graphicFileTypes, nil
}

func CreateGraphic(ctx context.Context, graphic *Graphic) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	graphic.Created = time.Now()

	key := datastore.IncompleteKey("Graphic", nil)
	putKey, err := client.Put(ctx, key, graphic)
	if err != nil {
		return err
	}

	graphic.ID = putKey.ID

	return nil
}

func DeleteGraphicsInTransaction(tx *datastore.Transaction, ids []int64) error {
	var graphicKeys []*datastore.Key

	for _, id := range ids {
		graphicKeys = append(graphicKeys, datastore.IDKey("Graphic", id, nil))
	}

	if err := tx.DeleteMulti(graphicKeys); err != nil {
		return err
	}

	return nil
}

func storeGraphicThumbnailURL(ctx context.Context, graphics *[]Graphic, keys []*datastore.Key) error {
	for i, key := range keys {
		fileID := strconv.FormatInt(key.ID, 10)
		filePath := storageObjectFilePath("Graphic", fileID, (*graphics)[i].FileType)
		fileType := "" // this is unnecessary when GET request
		bucketName := infrastructure.MaterialBucketName()
		url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)
		if err != nil {
			return err
		}

		(*graphics)[i].ID = key.ID
		(*graphics)[i].URL = url
		(*graphics)[i].ThumbnailURL = infrastructure.GraphicThumbnailURL(ctx, key.ID, fileType)
	}

	return nil
}
