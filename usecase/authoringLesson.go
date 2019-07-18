package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type AuthoringLessonErrorCode uint

const (
	AuthoringLessonNotFound      AuthoringLessonErrorCode = 1
	InvalidAuthoringLessonParams AuthoringLessonErrorCode = 2
)

func (e AuthoringLessonErrorCode) Error() string {
    switch e {
    case AuthoringLessonNotFound:
        return "authoring lesson not found"
    case InvalidAuthoringLessonParams:
        return "invalid authoring lesson params"
    default:
        return "unknown authoring lesson error"
    }
}

func GetAuthoringLesson(request *http.Request, id string) (domain.Lesson, error) {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := getLessonByIdWithResources(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, AuthoringLessonNotFound
	} else {
		return lesson, err
	}

	if lesson.UserID != currentUser.ID {
		return lesson, InvalidAuthoringLessonParams
	}

	return lesson, nil
}

func CreateAuthoringLesson(request *http.Request, lesson domain.Lesson) (domain.Lesson, error) {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return lesson, err
	}

	id := xid.New().String()
	lesson.ID = id
	lesson.UserID = currentUser.ID
	lesson.Created = time.Now()

	if err = domain.CreateNewLesson(ctx, lesson); err != nil {
		return lesson, err
	}

	return lesson, nil
}

func UpdateAuthoringLesson(id string, request *http.Request) (domain.Lesson, error){
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := domain.GetLessonById(ctx, id)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lesson, AuthoringLessonNotFound
		}
		return lesson, err
	}
	lesson.Updated = time.Now()

	if lesson.UserID != currentUser.ID {
		return lesson, InvalidAuthoringLessonParams
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, request.Body)

	var f interface{}
	if err := json.Unmarshal(buf.Bytes(), &f); err != nil {
		return lesson, InvalidAuthoringLessonParams
	}

	updateLesson := f.(map[string]interface{})
	mutable := reflect.ValueOf(lesson).Elem()
	for key, lessonField := range updateLesson {
		structKey := strings.Title(key)
		switch v := lessonField.(type) {
		case []interface{}:
			array := make([]string, len(v)) // TODO support another types. reflect.TypeOf(v[0])
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

	if err = domain.UpdateLesson(ctx, lesson); err != nil {
		return lesson, err
	}

	return lesson, nil
}

func getLessonByIdWithResources(ctx context.Context, id string) (domain.Lesson, error) {
	lesson, err := domain.GetLessonById(ctx, id)

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lesson, AuthoringLessonNotFound
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
