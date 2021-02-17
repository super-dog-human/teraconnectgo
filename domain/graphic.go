package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Graphic is used for lesson.
type Graphic struct {
	ID           int64     `json:"id" datastore:"-"`
	LessonID     int64     `json:"lessonID"`
	FileType     string    `json:"fileType"`
	IsPublic     bool      `json:"isPublic"`
	URL          string    `json:"url" datastore:"-"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	Created      time.Time `json:"created"`
}

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

func CreateGraphics(ctx context.Context, userID int64, graphics []*Graphic) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	parentKey := datastore.IDKey("User", userID, nil)

	keys := make([]*datastore.Key, len(graphics))
	currentTime := time.Now()
	for i, graphic := range graphics {
		keys[i] = datastore.IncompleteKey("Graphic", parentKey)
		graphic.Created = currentTime
	}

	putKeys, err := client.PutMulti(ctx, keys, graphics)
	if err != nil {
		return err
	}

	for i, graphic := range graphics {
		graphic.ID = putKeys[i].ID
	}

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
