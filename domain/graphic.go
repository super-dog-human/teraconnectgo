package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type GraphicErrorCode uint

const (
	GraphicNotFound GraphicErrorCode = 1
)

func (e GraphicErrorCode) Error() string {
	switch e {
	case GraphicNotFound:
		return "graphic not found"
	default:
		return "unknown graphic error"
	}
}

// Graphic is used for lesson.
type Graphic struct {
	ID       int64     `json:"id" datastore:"-"`
	LessonID int64     `json:"lessonID"`
	FileType string    `json:"fileType"`
	IsPublic bool      `json:"isPublic"`
	URL      string    `json:"url" datastore:"-"`
	Created  time.Time `json:"created"`
}

func GetGraphicByID(ctx context.Context, id int64, userID int64) (Graphic, error) {
	graphic := new(Graphic)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *graphic, err
	}

	ancestor := datastore.IDKey("User", userID, nil)
	key := datastore.IDKey("Graphic", id, ancestor)

	if err := client.Get(ctx, key, graphic); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *graphic, GraphicNotFound
		}
		return *graphic, err
	}

	graphic.ID = id

	return *graphic, nil
}

func GetGraphicsByLessonID(ctx context.Context, lessonID int64, graphics *[]Graphic) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	query := datastore.NewQuery("Graphic").Filter("LessonID =", lessonID).Order("Created")

	keys, err := client.GetAll(ctx, query, graphics)
	if err != nil {
		return err
	}

	if len(*graphics) == 0 {
		return GraphicNotFound
	}

	storeGraphicURL(ctx, graphics, keys)

	return nil
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

func DeleteGraphicByID(ctx context.Context, id int64, userID int64) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	ancestor := datastore.IDKey("User", userID, nil)
	key := datastore.IDKey("Graphic", id, ancestor)
	if err := client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}

func DeleteGraphicFileByID(ctx context.Context, graphic Graphic) error {
	bucketName := infrastructure.MaterialBucketName()
	fileID := strconv.FormatInt(graphic.ID, 10)
	filePath := infrastructure.StorageObjectFilePath("Graphic", fileID, graphic.FileType)

	if err := infrastructure.DeleteObjectFromGCS(ctx, bucketName, filePath); err != nil {
		return err
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

func storeGraphicURL(ctx context.Context, graphics *[]Graphic, keys []*datastore.Key) error {
	for i, key := range keys {
		fileID := strconv.FormatInt(key.ID, 10)
		filePath := infrastructure.StorageObjectFilePath("Graphic", fileID, (*graphics)[i].FileType)
		fileType := "" // this is unnecessary when GET request
		bucketName := infrastructure.MaterialBucketName()
		url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)
		if err != nil {
			return err
		}

		(*graphics)[i].ID = key.ID
		(*graphics)[i].URL = url
	}

	return nil
}
