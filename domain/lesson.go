package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/rs/xid"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetLessonByID(ctx context.Context, id string) (Lesson, error) {
	lesson := new(Lesson)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *lesson, err
	}

	key := datastore.NameKey("Lesson", id, nil)
	if err := client.Get(ctx, key, lesson); err != nil {
		return *lesson, err
	}
	lesson.ID = id

	return *lesson, nil
}

func GetLessonsByUserID(ctx context.Context, userID string) ([]Lesson, error) {
	var lessons []Lesson

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("Lesson").Filter("UserID =", userID)
	if _, err := client.GetAll(ctx, query, &lessons); err != nil {
		return nil, err
	}

	return lessons, nil
}

func CreateLesson(ctx context.Context, lesson *Lesson, userID string) error {
	lesson.ID = xid.New().String()
	lesson.UserID = userID
	lesson.Created = time.Now()

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("Lesson", lesson.ID, nil)
	if _, err := client.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func UpdateLesson(ctx context.Context, lesson *Lesson) error {
	lesson.Updated = time.Now()

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("Lesson", lesson.ID, nil)
	if _, err := client.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func DeleteLessonInTransactionByID(tx *datastore.Transaction, id string) error {
	key := datastore.NameKey("Lesson", id, nil)
	if err := tx.Delete(key); err != nil {
		return err
	}

	return nil
}
