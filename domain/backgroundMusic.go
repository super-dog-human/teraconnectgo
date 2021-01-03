package domain

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// BackgroundMusic type is used in the class.
type BackgroundMusic struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	SortID int64  `json:"-"`
}

// GetAllBackgroundMusics is return all sorted musics.
func GetAllBackgroundMusics(ctx context.Context) ([]BackgroundMusic, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var musics []BackgroundMusic
	query := datastore.NewQuery("BackgroundMusic").Order("SortID")
	keys, err := client.GetAll(ctx, query, &musics)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		musics[i].ID = key.ID
		musics[i].URL = infrastructure.GetPublicBackGroundMusicURL(strconv.FormatInt(key.ID, 10))
	}

	return musics, nil
}
