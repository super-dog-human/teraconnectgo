package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type LessonErrorCode uint

const (
	LessonNotFound     LessonErrorCode = 1
	LessonNotAvailable LessonErrorCode = 2
)

func (e LessonErrorCode) Error() string {
    switch e {
	case LessonNotFound:
		return "lesson not found"
	case LessonNotAvailable:
		return "lesson not available"
    default:
        return "unknown lesson error"
    }
}

// GetLessons for fetch lessons
func GetLessons(request *http.Request) ([]domain.Lesson, error) {
	// 検索パラメータ
	// SearchAPIが必須
	return nil, nil
}

// GetLesson for fetch the lesson by id
func GetAvailableLesson(request *http.Request, id string) (domain.Lesson, error) {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	authErr, ok := err.(domain.AuthErrorCode)
	if !ok || authErr != domain.TokenNotFound {
		// can get lesson without token, but can NOT get with invalid token.
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := getLessonById(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, LessonNotFound
	} else {
		return lesson, err
	}

	if lesson.IsPublic {
		return lesson, nil
	}

	fmt.Printf("currentUser ID %v+\n", currentUser.ID)
	if currentUser.ID != "" && lesson.UserID == currentUser.ID {
		return lesson, nil
	}

	return lesson, LessonNotAvailable
}

func DestroyOwnLessonById(request *http.Request, id string) error {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := getLessonById(ctx, id)
	if err != nil{
		if err == datastore.ErrNoSuchEntity {
			return LessonNotFound
		}
		return err
	}

	if currentUser.ID != lesson.UserID {
		return LessonNotAvailable
	}


	if err := destroyLessonAndRecources(ctx, lesson.ID); err != nil {
		return err
	}

	return nil
}

func getLessonById(ctx context.Context, id string) (domain.Lesson, error) {
	lesson := new(domain.Lesson)

	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if err := datastore.Get(ctx, key, lesson); err != nil {
		return *lesson, err
	}
	lesson.ID = id

	return *lesson, nil
}

func createNewLesson(ctx context.Context, lesson domain.Lesson) error {
	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func updateLessonById(ctx context.Context, lesson domain.Lesson) error {
	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

func destroyLessonAndRecources(ctx context.Context, id string) error {
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
