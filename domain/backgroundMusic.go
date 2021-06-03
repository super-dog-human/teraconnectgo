package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// BackgroundMusic type is used in the class.
type BackgroundMusic struct {
	ID       int64     `json:"id" datastore:"-"`
	Name     string    `json:"name"`
	URL      string    `json:"url" datastore:"-"`
	SortID   int64     `json:"-"`
	IsPublic bool      `json:"-"`
	Created  time.Time `json:"created"`
}

// GetPublicBackgroundMusics is return sorted public musics.
func GetPublicBackgroundMusics(ctx context.Context) ([]BackgroundMusic, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var musics []BackgroundMusic
	query := datastore.NewQuery("BackgroundMusic").Filter("IsPublic =", true).Order("SortID")
	keys, err := client.GetAll(ctx, query, &musics)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		musics[i].ID = key.ID
		musics[i].URL = infrastructure.GetPublicBackgroundMusicURL(strconv.FormatInt(key.ID, 10))
	}

	return musics, nil
}

func GetCurrentUsersBackgroundMusics(ctx context.Context, userID int64) ([]BackgroundMusic, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var musics []BackgroundMusic
	ancestor := datastore.IDKey("User", userID, nil)
	query := datastore.NewQuery("BackgroundMusic").Ancestor(ancestor).Order("-Created")
	keys, err := client.GetAll(ctx, query, &musics)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		musics[i].ID = key.ID
		url, err := getBackgroundMusicSignedURL(ctx, key.ID)
		if err != nil {
			return nil, err
		}
		musics[i].URL = url
	}

	return musics, nil
}

func CreateBackgroundMusic(ctx context.Context, userID int64, backgroundMusic *BackgroundMusic) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	backgroundMusic.Created = time.Now()

	ancestor := datastore.IDKey("User", userID, nil)
	key := datastore.IncompleteKey("BackgroundMusic", ancestor)
	putKey, err := client.Put(ctx, key, backgroundMusic)
	if err != nil {
		return err
	}

	backgroundMusic.ID = putKey.ID

	return nil
}

func getBackgroundMusicSignedURL(ctx context.Context, id int64) (string, error) {
	fileID := strconv.FormatInt(id, 10)
	filePath := infrastructure.StorageObjectFilePath("bgm", fileID, "mp3")
	fileType := "" // this is unnecessary when GET request
	bucketName := infrastructure.MaterialBucketName()
	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)

	if err != nil {
		return "", err
	}

	return url, nil
}
