package domain

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Subject of the class type.
type Subject struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	JapaneseName string `json:"japaneseName"`
	SortID       int64  `json:"-"`
}

// GetAllSubjects is return all sorted subjects.
func GetAllSubjects(ctx context.Context) ([]Subject, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	var subjects []Subject
	query := datastore.NewQuery("Subject").Order("SortID")
	keys, err := client.GetAll(ctx, query, &subjects)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		subjects[i].ID = key.ID
	}

	return subjects, nil
}
