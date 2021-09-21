package domain

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type CategoryErrorCode uint

const (
	CategoryNotFound CategoryErrorCode = 1
)

func (e CategoryErrorCode) Error() string {
	switch e {
	case CategoryNotFound:
		return "category not found"
	default:
		return "unknown category error"
	}
}

// Category of the class type.
type Category struct {
	ID          int64  `json:"id" datastore:"-"`
	SubjectID   int64  `json:"subjectID"`
	SubjectName string `json:"subjectName" datastore:",noindex"`
	GroupName   string `json:"groupName" datastore:",noindex"`
	Name        string `json:"name" datastore:",noindex"`
	SortID      int64  `json:"-"`
}

type ShortCategory struct {
	ID          int64  `json:"id"`
	SubjectID   int64  `json:"-"`
	SubjectName string `json:"-"`
	GroupName   string `json:"groupName"`
	Name        string `json:"name"`
	SortID      int64  `json:"-"`
}

// GetJapaneseCategories is return categories by the subject.
func GetJapaneseCategories(ctx context.Context, subjectID int64) ([]ShortCategory, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var categories []ShortCategory
	query := datastore.NewQuery("JapaneseCategory").Filter("SubjectID =", subjectID).Order("SortID")
	keys, err := client.GetAll(ctx, query, &categories)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		categories[i].ID = key.ID
	}

	return categories, nil
}

// GetAllJapaneseCategories is return all sorted categories.
func GetAllJapaneseCategories(ctx context.Context) ([]ShortCategory, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var categories []ShortCategory
	query := datastore.NewQuery("JapaneseCategory").Order("SubjectID").Order("SortID")
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

	key := datastore.IDKey("JapaneseCategory", id, nil)
	if err := client.Get(ctx, key, category); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *category, CategoryNotFound
		}
		return *category, err
	}
	category.ID = id

	return *category, nil
}
