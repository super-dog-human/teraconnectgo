package usecase

import (
	"log"
	"net/http"
	"reflect"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
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
	SubjectID          int64  `json:"subjectID"`
	JapaneseCategoryID int64  `json:"japaneseCategoryID"`
	Title              string `json:"title"`
}

type PatchLessonParams struct {
	SubjectID          int64                     `json:"subjectID"`
	JapaneseCategoryID int64                     `json:"japaneseCategoryID"`
	Status             domain.LessonStatus       `json:"status"`
	Title              string                    `json:"title"`
	Description        string                    `json:"description"`
	References         []domain.LessonReferences `json:"feferences"`
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

	if lesson.Status == domain.LessonStatusLimited {
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

	subject, err := domain.GetSubject(ctx, newLesson.SubjectID)
	if err != nil {
		return InvalidLessonParams
	}

	category, err := domain.GetJapaneseCategory(ctx, newLesson.JapaneseCategoryID, newLesson.SubjectID)
	if err != nil {
		log.Printf("category %v\n", newLesson.JapaneseCategoryID)
		log.Printf("%v\n", errors.WithStack(err).Error())
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

func UpdateLessonWithMaterial(id int64, materialID int64, request *http.Request, lessonParams *PatchLessonParams, materialParams *PatchLessonMaterialParams) error {
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

	var newLesson domain.Lesson
	var newLessonMaterial domain.LessonMaterial

	var blankLessonParams PatchLessonParams
	if !reflect.DeepEqual(*lessonParams, blankLessonParams) {
		copier.Copy(&newLesson, *lessonParams)
		newLesson.ID = id
	}

	var blankMaterialParams PatchLessonParams
	if !reflect.DeepEqual(*materialParams, blankMaterialParams) {
		copier.Copy(&newLessonMaterial, *materialParams)
		newLessonMaterial.ID = materialID
	}

	if err := domain.UpdateLessonAndMaterial(ctx, &newLesson, &newLessonMaterial); err != nil {
		return err
	}

	return nil
}
