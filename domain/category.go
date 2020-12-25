package domain

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Category of the class type.
type Category struct {
	ID        int64  `json:"id"`
	GroupName string `json:"groupName"`
	Name      string `json:"name"`
	SortID    int64  `json:"-"`
}

// GetJapaneseCategories is return categories by the subject.
func GetJapaneseCategories(ctx context.Context, subjectID int64) ([]Category, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var categories []Category
	ancestor := datastore.IDKey("Subject", subjectID, nil)
	query := datastore.NewQuery("JapaneseCategory").Ancestor(ancestor).Order("SortID")
	keys, err := client.GetAll(ctx, query, &categories)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		categories[i].ID = key.ID
	}

	return categories, nil
}

// GetJapaneseCategory is return a category from id.
func GetJapaneseCategory(ctx context.Context, id int64, subjectID int64) (Category, error) {
	category := new(Category)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *category, err
	}

	ancestor := datastore.IDKey("Subject", subjectID, nil)
	key := datastore.IDKey("JapaneseCategory", id, ancestor)
	if err := client.Get(ctx, key, category); err != nil {
		return *category, err
	}
	category.ID = id

	return *category, nil
}
