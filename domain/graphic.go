package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetGraphicsByIDs(ctx context.Context, ids []string) ([]Graphic, error) {
	var graphicKeys []*datastore.Key

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		graphicKeys = append(graphicKeys, datastore.NameKey("Graphic", id, nil))
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

func GetCurrentUsersGraphics(ctx context.Context, userID string) ([]Graphic, error) {
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

func GetGraphicFileTypes(ctx context.Context, graphicIDs []string) (map[string]string, error) {
	var keys []*datastore.Key

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	for _, id := range graphicIDs {
		keys = append(keys, datastore.NameKey("Graphic", id, nil))
	}

	graphicFileTypes := map[string]string{}
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

func CreateGraphic(ctx context.Context, id string, userID string, fileType string) error {
	graphic := new(Graphic)

	graphic.ID = id
	graphic.UserID = userID
	graphic.Created = time.Now()
	graphic.FileType = fileType

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("Graphic", graphic.ID, nil)
	if _, err := client.Put(ctx, key, graphic); err != nil {
		return err
	}

	return nil
}

func DeleteGraphicsInTransaction(tx *datastore.Transaction, ids []string) error {
	var graphicKeys []*datastore.Key

	for _, id := range ids {
		graphicKeys = append(graphicKeys, datastore.NameKey("Graphic", id, nil))
	}

	if err := tx.DeleteMulti(graphicKeys); err != nil {
		return err
	}

	return nil
}

func storeGraphicThumbnailURL(ctx context.Context, graphics *[]Graphic, keys []*datastore.Key) error {
	for i, key := range keys {
		id := key.Name
		filePath := storageObjectFilePath("Graphic", id, (*graphics)[i].FileType)
		fileType := "" // this is unnecessary when GET request
		bucketName := infrastructure.MaterialBucketName()
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
