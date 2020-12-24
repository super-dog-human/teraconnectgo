package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetLessonByID(ctx context.Context, id int64) (Lesson, error) {
	lesson := new(Lesson)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *lesson, err
	}

	key := datastore.IDKey("Lesson", id, nil)
	if err := client.Get(ctx, key, lesson); err != nil {
		return *lesson, err
	}
	lesson.ID = id

	return *lesson, nil
}

func GetLessonsByUserID(ctx context.Context, userID int64) ([]Lesson, error) {
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

func CreateLesson(ctx context.Context, lesson *Lesson) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	lesson.Created = time.Now()

	key := datastore.IncompleteKey("Lesson", nil)
	puttedKey, err := client.Put(ctx, key, lesson)

	if err != nil {
		return err
	}

	lesson.ID = puttedKey.ID

	return nil
}

func UpdateLesson(ctx context.Context, lesson *Lesson) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	lesson.Updated = time.Now()

	key := datastore.IDKey("Lesson", lesson.ID, nil)
	if _, err := client.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func DeleteLessonInTransactionByID(tx *datastore.Transaction, id int64) error {
	key := datastore.IDKey("Lesson", id, nil)
	if err := tx.Delete(key); err != nil {
		return err
	}

	return nil
}
