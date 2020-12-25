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

// GetSubject is return a subject from id.
func GetSubject(ctx context.Context, id int64) (Subject, error) {
	subject := new(Subject)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *subject, err
	}

	key := datastore.IDKey("Subject", id, nil)
	if err := client.Get(ctx, key, subject); err != nil {
		return *subject, err
	}
	subject.ID = id

	return *subject, nil
}
