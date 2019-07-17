package domain

import (
	"context"

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

func CreateNewLesson(ctx context.Context, lesson Lesson) error {
	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func UpdateLesson(ctx context.Context, lesson Lesson) error {
	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func DestroyLessonAndRecources(ctx context.Context, id string) error {
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := destroyLessonById(ctx, id); err != nil {
			return err
		}

		// remove voicetexts
		// remove voice files
		// remove zip file

//		_, err := datastore.Delete(ctx, lessonKey, lesson)
		return nil
	}, nil)

	return err
}

func destroyLessonById(ctx context.Context, id string) error {
	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if err := datastore.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
