package domain

import (
	"context"
	"time"

	"github.com/rs/xid"
	"google.golang.org/appengine/datastore"
)

func GetLessonById(ctx context.Context, id string) (Lesson, error) {
	lesson := new(Lesson)

	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if err := datastore.Get(ctx, key, lesson); err != nil {
		return *lesson, err
	}
	lesson.ID = id

	return *lesson, nil
}

func GetLessonsByUserID(ctx context.Context, userID string) ([]Lesson, error) {
	var lessons []Lesson

	query := datastore.NewQuery("Lesson").Filter("UserID =", userID)
	if _, err := query.GetAll(ctx, &lessons); err != nil {
		return lessons, err
	}

	return lessons, nil
}

func CreateNewLesson(ctx context.Context, lesson Lesson, userID string) error {
	lesson.ID = xid.New().String()
	lesson.UserID = userID
	lesson.Created = time.Now()

	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func UpdateLesson(ctx context.Context, lesson Lesson) error {
	lesson.Updated = time.Now()
	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func DeleteLessonById(ctx context.Context, id string) error {
	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if err := datastore.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
