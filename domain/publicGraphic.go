package domain

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// PublicGraphic is image entity used for introduction.
type PublicGraphic struct {
	ID              int64  `datastore:"-"`
	FileType        string `datastore:",noindex"`
	ForIntroduction bool
	SortID          int64
}

func GetPublicGraphicsForIntroduction(ctx context.Context) ([]PublicGraphic, error) {
	var graphics []PublicGraphic

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("PublicGraphic").Filter("ForIntroduction =", true).Order("SortID")

	keys, err := client.GetAll(ctx, query, &graphics)
	if err != nil {
		return graphics, err
	}

	for i, key := range keys {
		graphics[i].ID = key.ID
	}

	return graphics, nil
}
