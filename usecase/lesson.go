package usecase

import (
	"context"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/domain"
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

// NewLessonParamsは、Lessonの新規作成時、リクエストボディをbindするために使用されます。
type NewLessonParams struct {
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

	if err = setRelationLessonTitle(ctx, &lesson); err != nil {
		return lesson, err
	}

	if lesson.Status == domain.LessonStatusPublic {
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

	ctx := request.Context()
	lesson, err := domain.GetLessonByID(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, LessonNotFound
	} else if err != nil {
		return lesson, err
	}

	if lesson.UserID != currentUser.ID {
		return lesson, InvalidLessonParams
	}

	if err = setRelationLessonTitle(ctx, &lesson); err != nil {
		return lesson, err
	}

	return lesson, nil
}

func GetCurrentUserLessons(request *http.Request) ([]domain.Lesson, error) {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	lessons, err := domain.GetLessonsByUserID(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

// CreateLesson is create the new lesson belongs to subject and category.
func CreateLesson(request *http.Request, newLesson *NewLessonParams, lesson *domain.Lesson) error {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return InvalidLessonParams
	}

	copier.Copy(&lesson, &newLesson)

	ctx := request.Context()

	lesson.UserID = currentUser.ID

	if err = domain.CreateLesson(ctx, lesson); err != nil {
		return err
	}

	return nil
}

func UpdateLessonWithMaterial(id int64, request *http.Request, params *map[string]interface{}) error {
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

	if lesson.UserID != currentUser.ID {
		return InvalidLessonParams
	}

	lessonFields := []string{"PrevLessonID", "NextLessonID", "SubjectID", "JapaneseCategoryID", "Status", "HasThumbnail", "Title", "Description", "References"}
	lessonMaterialFields := []string{"BackgroundImageID", "AvatarID", "AvatarLightColor", "VoiceSynthesisConfig"}
	if err := domain.UpdateLessonAndMaterial(ctx, &lesson, params, &lessonFields, &lessonMaterialFields); err != nil {
		return err
	}

	return nil
}

func setRelationLessonTitle(ctx context.Context, lesson *domain.Lesson) error {
	if lesson.PrevLessonID != 0 {
		lesson, err := domain.GetLessonByID(ctx, lesson.PrevLessonID)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return nil // 授業が見つからなかった場合もエラーにしない
			}
			return err
		}
		lesson.PrevLessonTitle = lesson.Title
	}

	if lesson.NextLessonID != 0 {
		lesson, err := domain.GetLessonByID(ctx, lesson.NextLessonID)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return nil // 授業が見つからなかった場合もエラーにしない
			}
			return err
		}
		lesson.NextLessonTitle = lesson.Title
	}

	return nil
}
