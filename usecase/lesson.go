package usecase

import (
	"context"
	"net/http"
	"reflect"

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

// NewLessonParams is only used when creating new lesson.
type NewLessonParams struct {
	NeedsRecording     bool   `json:"needsRecording"`
	IsIntroduction     bool   `json:"isIntroduction"`
	HasThumbnail       bool   `json:"hasThumbnail"`
	SubjectID          int64  `json:"subjectID"`
	JapaneseCategoryID int64  `json:"japaneseCategoryID"`
	Title              string `json:"title"`
}

type PatchLessonAndMaterialParams struct {
	PatchLessonParams
	PatchLessonMaterialParams
}

type PatchLessonParams struct {
	PrevLessonID       int64                     `json:"prevLessonID"`
	NextLessonID       int64                     `json:"nextLessonID"`
	SubjectID          int64                     `json:"subjectID"`
	JapaneseCategoryID int64                     `json:"japaneseCategoryID"`
	Status             domain.LessonStatus       `json:"status"`
	Title              string                    `json:"title"`
	Description        string                    `json:"description"`
	References         []domain.LessonReferences `json:"references"`
}

type PatchLessonMaterialParams struct {
	BackgroundImageID    int64                       `json:"backgroundImageID"`
	AvatarID             int64                       `json:"avatarID"`
	AvatarLightColor     string                      `json:"avatarLightColor"`
	VoiceSynthesisConfig domain.VoiceSynthesisConfig `json:"voiceSynthesisConfig"`
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

	if err := setCategoryAndSubject(ctx, newLesson.SubjectID, newLesson.JapaneseCategoryID, lesson); err != nil {
		return InvalidLessonParams
	}

	lesson.UserID = currentUser.ID

	if err = domain.CreateLesson(ctx, lesson); err != nil {
		return err
	}

	return nil
}

func UpdateLessonWithMaterial(id int64, request *http.Request, params *PatchLessonAndMaterialParams) error {
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

	var lessonParams PatchLessonParams
	var materialParams PatchLessonMaterialParams

	copier.Copy(&lessonParams, *params)
	copier.Copy(&materialParams, *params)

	var newLesson domain.Lesson
	var newLessonMaterial domain.LessonMaterial

	var blankLessonParams PatchLessonParams
	if !reflect.DeepEqual(lessonParams, blankLessonParams) {
		copier.Copy(&newLesson, lessonParams)
	}

	if lessonParams.SubjectID != 0 && lessonParams.JapaneseCategoryID != 0 {
		if err = setCategoryAndSubject(ctx, lessonParams.SubjectID, lessonParams.JapaneseCategoryID, &lesson); err != nil {
			return err
		}
	}

	var blankMaterialParams PatchLessonParams
	if !reflect.DeepEqual(materialParams, blankMaterialParams) {
		copier.Copy(&newLessonMaterial, materialParams)
	}

	if err := domain.UpdateLessonAndMaterial(ctx, id, lesson.MaterialID, &newLesson, &newLessonMaterial); err != nil {
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

func setCategoryAndSubject(ctx context.Context, subjectID int64, japaneseCategoryID int64, lesson *domain.Lesson) error {
	subject, err := domain.GetSubject(ctx, subjectID)
	if err != nil {
		return err
	}

	category, err := domain.GetJapaneseCategory(ctx, japaneseCategoryID, subjectID)
	if err != nil {
		return err
	}

	lesson.SubjectName = subject.JapaneseName
	lesson.JapaneseCategoryName = category.Name

	return nil
}
