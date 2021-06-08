package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
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

// LessonParams is only used when creating new lesson.
type LessonParams struct {
	NeedsRecording     bool   `json:"needsRecording"`
	IsIntroduction     bool   `json:"isIntroduction"`
	SubjectID          int64  `json:"subjectID"`
	JapaneseCategoryID int64  `json:"japaneseCategoryID"`
	Title              string `json:"title"`
}

// GetLessonsByConditions for search lessons
func GetLessonsByConditions(request *http.Request) ([]domain.Lesson, error) {
	//		ctx := request.Context()

	//	id := c.Param("id")

	// 検索パラメータ
	// ページネーション
	// SearchAPIが必須

	return nil, nil
}

// GetPublicLesson for fetch the lesson by id
func GetPublicLesson(request *http.Request, id int64) (domain.Lesson, error) {
	ctx := request.Context()

	_, err := domain.GetCurrentUser(request)
	authErr, _ := err.(domain.AuthErrorCode)

	if err != nil && authErr != domain.TokenNotFound {
		// can get lesson without token, but can NOT get with invalid token.
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := domain.GetLessonByID(ctx, id)
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

func GetPrivateLesson(request *http.Request, id int64) (domain.Lesson, error) {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := getLessonByIDWithResources(request, id)
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

// CreateLesson is create the new lesson belongs to subject and category.
func CreateLesson(request *http.Request, newLesson *LessonParams, lesson *domain.Lesson) error {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return InvalidLessonParams
	}

	copier.Copy(&lesson, &newLesson)

	ctx := request.Context()

	subject, err := domain.GetSubject(ctx, newLesson.SubjectID)
	if err != nil {
		return InvalidLessonParams
	}

	category, err := domain.GetJapaneseCategory(ctx, newLesson.JapaneseCategoryID, newLesson.SubjectID)
	if err != nil {
		log.Printf("category %v\n", newLesson.JapaneseCategoryID)
		log.Printf("%v\n", errors.WithStack(err.(error)).Error())
		return InvalidLessonParams
	}

	lesson.UserID = currentUser.ID
	lesson.SubjectName = subject.JapaneseName
	lesson.CategoryName = category.Name

	if err = domain.CreateLesson(ctx, lesson); err != nil {
		return err
	}

	return nil
}

func UpdateLesson(id int64, request *http.Request) (domain.Lesson, error) {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		lesson := new(domain.Lesson)
		return *lesson, err
	}

	lesson, err := domain.GetLessonByID(ctx, id)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lesson, LessonNotFound
		}
		return lesson, err
	}

	// TODO allow permitted users for authoring
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

func DeleteOwnLessonByID(request *http.Request, id int64) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLessonByID(ctx, id)
	if err != nil {
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

func getLessonByIDWithResources(request *http.Request, id int64) (domain.Lesson, error) {
	var lesson domain.Lesson

	ctx := request.Context()
	lesson, err := domain.GetLessonByID(ctx, id)

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lesson, LessonNotFound
		}
		return lesson, err
	}

	if lesson.AvatarID != 0 {
		avatar, err := domain.GetPublicAvatarByID(ctx, lesson.AvatarID)
		if err != nil {
			if ok := errors.Is(err, domain.AvatarNotFound); ok {
				currentUser, err := domain.GetCurrentUser(request)
				if err != nil {
					return lesson, err
				}

				avatar, err = domain.GetCurrentUsersAvatarByID(ctx, lesson.AvatarID, currentUser.ID)
				if err != nil {
					return lesson, err
				}
			} else {
				return lesson, err
			}
		}

		lesson.Avatar = avatar
	}

	return lesson, nil
}

func deleteLessonAndRecources(ctx context.Context, lesson domain.Lesson) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := domain.DeleteAvatarInTransaction(tx, lesson.AvatarID); err != nil {
			return err
		}

		if err := domain.DeleteGraphicsInTransaction(tx, lesson.GraphicIDs); err != nil {
			return err
		}

		// TODO remove files in GCS

		if err := domain.DeleteLessonInTransactionByID(tx, lesson.ID); err != nil {
			return err
		}

		return nil
	})

	return err
}
