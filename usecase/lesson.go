package usecase

import (
	"context"
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

	lesson, err := domain.GetLessonById(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, LessonNotFound
	} else {
		return lesson, err
	}

	if lesson.IsPublic {
		return lesson, nil
	}

	if currentUser.ID != "" && lesson.UserID == currentUser.ID {
		return lesson, nil
	}

	return lesson, LessonNotAvailable
}

func DeleteOwnLessonById(request *http.Request, id string) error {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLessonById(ctx, id)
	if err != nil{
		if err == datastore.ErrNoSuchEntity {
			return LessonNotFound
		}
		return err
	}

	if currentUser.ID != lesson.UserID {
		return LessonNotAvailable
	}

	if err := deleteLessonAndRecources(ctx, lesson); err != nil {
		return err
	}

	return nil
}

func deleteLessonAndRecources(ctx context.Context, lesson domain.Lesson) error {
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {

		if err := domain.DeleteAvatar(ctx, lesson.AvatarID); err != nil {
			return err
		}

		if err := domain.DeleteGraphics(ctx, lesson.GraphicIDs); err != nil {
			return err
		}

		if err := domain.DeleteRawVoiceTextsByLessonID(ctx, lesson.ID); err != nil {
			return err
		}

		// TODO remove files in GCS

		if err := domain.DeleteLessonById(ctx, lesson.ID); err != nil {
			return err
		}

		return nil
	}, nil)

	return err
}
