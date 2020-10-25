package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/super-dog-human/teraconnectgo/domain"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type LessonErrorCode uint

const (
	LessonNotFound      LessonErrorCode = 1
	LessonNotAvailable  LessonErrorCode = 2
	InvalidLessonParams LessonErrorCode = 3
)

func (e LessonErrorCode) Error() string {
    switch e {
	case LessonNotFound:
		return "lesson not found"
	case LessonNotAvailable:
		return "lesson not available"
	case InvalidLessonParams:
        return "invalid lesson params"
    default:
        return "unknown lesson error"
    }
}

// GetLessonsByConditions for search lessons
func GetLessonsByConditions(request *http.Request) ([]domain.Lesson, error) {
//		ctx := appengine.NewContext(request)

//	id := c.Param("id")

	// 検索パラメータ
	// ページネーション
	// SearchAPIが必須

	return nil, nil
}

// GetPublicLesson for fetch the lesson by id
func GetPublicLesson(request *http.Request, id string) (domain.Lesson, error) {
	ctx := appengine.NewContext(request)

	_, err := domain.GetCurrentUser(request)
	authErr, ok := err.(domain.AuthErrorCode)
	if !ok || authErr != domain.TokenNotFound {
		// can get lesson without token, but can NOT get with invalid token.
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := domain.GetLessonById(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, LessonNotFound
	} else if err != nil {
		return lesson, err
	}

	if lesson.IsPublic {
		return lesson, nil
	}

	return lesson, LessonNotAvailable
}

func GetPrivateLesson(request *http.Request, id string) (domain.Lesson, error) {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := getLessonByIdWithResources(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, LessonNotFound
	} else if err != nil {
		return lesson, err
	}

	if lesson.UserID != currentUser.ID {
		return lesson, InvalidLessonParams
	}

	return lesson, nil
}

func CreateLesson(request *http.Request, lesson *domain.Lesson) error {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	if err = domain.CreateLesson(ctx, lesson, currentUser.ID); err != nil {
		return err
	}

	// TODO upload thumbnail to GCS

	return nil
}

func UpdateLesson(id string, request *http.Request) (domain.Lesson, error){
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := domain.GetLessonById(ctx, id)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lesson, LessonNotFound
		}
		return lesson, err
	}

	if lesson.UserID != currentUser.ID {
		return lesson, InvalidLessonParams
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, request.Body)

	var f interface{}
	if err := json.Unmarshal(buf.Bytes(), &f); err != nil {
		return lesson, InvalidLessonParams
	}

	updateLesson := f.(map[string]interface{})
	mutable := reflect.ValueOf(lesson).Elem()
	for key, lessonField := range updateLesson {
		structKey := strings.Title(key)
		switch v := lessonField.(type) {
		case []interface{}:
			array := make([]string, len(v)) // TODO support not string in array types. use reflect.TypeOf(v[0])
			mutable.FieldByName(structKey).Set(reflect.ValueOf(array))
			for i := range v {
				mutable.FieldByName(structKey).Index(i).Set(reflect.ValueOf(v[i]))
			}
		default:
			if structKey == "ViewCount" || structKey == "Version" {
				intValue := int64(v.(float64))
				mutable.FieldByName(structKey).SetInt(intValue)
			} else {
				mutable.FieldByName(structKey).Set(reflect.ValueOf(v))
			}
		}
	}

	if err = domain.UpdateLesson(ctx, &lesson); err != nil {
		return lesson, err
	}

	return lesson, nil
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

func getLessonByIdWithResources(ctx context.Context, id string) (domain.Lesson, error) {
	lesson, err := domain.GetLessonById(ctx, id)

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lesson, LessonNotFound
		}
		return lesson, err
	}

	avatar, err := domain.GetAvatarByIds(ctx, lesson.AvatarID)
	if err != nil {
		return lesson, err
	}
	lesson.Avatar = avatar

	graphics, err := domain.GetGraphicsByIds(ctx, lesson.GraphicIDs)
	if err != nil {
		return lesson, err
	}
	lesson.Graphics = graphics

	return lesson, nil
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
