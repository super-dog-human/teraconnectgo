package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
	"golang.org/x/sync/errgroup"
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

// GetLessonsByCategoryID for search lessons
func GetLessonsByCategoryID(request *http.Request, categoryID int64, cursorStr string) ([]domain.ShortLesson, string, error) {
	ctx := request.Context()
	lessons, nextCursorStr, err := domain.GetLessonsByCategoryID(ctx, cursorStr, categoryID)
	if err != nil {
		return nil, "", err
	}
	return lessons, nextCursorStr, nil
}

// GetPublicLesson for fetch the lesson by id
func GetPublicLesson(request *http.Request, id int64, viewKey string) (domain.Lesson, error) {
	ctx := request.Context()

	lesson, err := domain.GetLessonByID(ctx, id)
	if err == datastore.ErrNoSuchEntity {
		return lesson, LessonNotFound
	} else if err != nil {
		return lesson, err
	}

	if lesson.Status != domain.LessonStatusPublic && lesson.Status != domain.LessonStatusLimited {
		return lesson, LessonNotAvailable
	}

	if lesson.Status == domain.LessonStatusLimited && lesson.ViewKey != viewKey {
		return lesson, LessonNotAvailable
	}

	if err = setRelationLessonTitle(ctx, &lesson); err != nil {
		return lesson, LessonNotAvailable
	}

	if err = setAvatar(ctx, &lesson); err != nil {
		return lesson, err
	}

	if err = setResourceURLs(ctx, &lesson); err != nil {
		return lesson, err
	}

	author, err := domain.GetUserByID(ctx, lesson.UserID)
	if err != nil {
		return lesson, nil
	}
	lesson.Author = author

	return lesson, nil
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

	lesson.Author = currentUser

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

	materialID, err := createInitialLessonMaterial(ctx, currentUser.ID, lesson.ID)
	if err != nil {
		return err
	}

	lesson.MaterialID = materialID

	if err := domain.UpdateLesson(ctx, lesson); err != nil {
		return err
	}

	return nil
}

// CreateIntroductionLessonは自己紹介用の授業を作成します。自己紹介に必要なGraphicも、初期データから作成します。
// 自己紹介授業は複数作成することはできないため、Lesson作成後にエラーが発生した場合はLessonの削除を試みます。
// Graphicはユーザーによる削除が可能なので、重複制限を行わず、Graphic作成後にエラーが発生してもロールバックは試みません。
// LessonMaterialも、Lesson削除後に残り続けても実害はないのでロールバックは試みません。
func CreateIntroductionLesson(request *http.Request, needsRecording bool, lesson *domain.Lesson) error {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return InvalidLessonParams
	}

	ctx := request.Context()

	lesson.NeedsRecording = needsRecording

	if err = domain.CreateIntroductionLesson(ctx, &currentUser, lesson); err != nil {
		return err
	}

	if err = domain.CreateIntroductionGraphics(ctx, currentUser.ID, lesson.ID); err != nil {
		// エラー時はLessonを削除する。削除時のエラーは無視する。
		domain.DeleteLesson(ctx, lesson.ID)
		return err
	}

	materialID, err := createInitialLessonMaterial(ctx, currentUser.ID, lesson.ID)
	if err != nil {
		// 同上
		domain.DeleteLesson(ctx, lesson.ID)
		return err
	}

	lesson.MaterialID = materialID

	if err := domain.UpdateLesson(ctx, lesson); err != nil {
		// 同上
		domain.DeleteLesson(ctx, lesson.ID)
		return err
	}

	return nil
}

func UpdateLessonWithMaterial(id int64, request *http.Request, needsCopyThumbnail bool, requestID string, params *map[string]interface{}) error {
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
	if err := domain.UpdateLessonAndMaterial(ctx, &lesson, needsCopyThumbnail, requestID, params, &lessonFields, &lessonMaterialFields); err != nil {
		return err
	}

	return nil
}

func setRelationLessonTitle(ctx context.Context, lesson *domain.Lesson) error {
	if lesson.PrevLessonID != 0 {
		prevLesson, err := domain.GetLessonByID(ctx, lesson.PrevLessonID)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return nil // 授業が見つからなかった場合もエラーにしない
			}
			return err
		}
		lesson.PrevLessonTitle = prevLesson.Title
	}

	if lesson.NextLessonID != 0 {
		nextLesson, err := domain.GetLessonByID(ctx, lesson.NextLessonID)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return nil // 授業が見つからなかった場合もエラーにしない
			}
			return err
		}
		lesson.NextLessonTitle = nextLesson.Title
	}

	return nil
}

func setAvatar(ctx context.Context, lesson *domain.Lesson) error {
	if lesson.AvatarID == 0 {
		return nil
	}

	avatar, err := domain.GetPublicAvatarByID(ctx, lesson.AvatarID)
	if err != nil {
		if ok := errors.Is(err, domain.AvatarNotFound); ok {
			avatar, err = domain.GetCurrentUsersAvatarByID(ctx, lesson.AvatarID, lesson.UserID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	lesson.Avatar = avatar

	return nil
}

func setResourceURLs(ctx context.Context, lesson *domain.Lesson) error {
	speechFilePath := fmt.Sprintf("lesson/%d/speech-%d.mp3", lesson.ID, lesson.Version)
	bodyFilePath := fmt.Sprintf("lesson/%d/body-%d.zst", lesson.ID, lesson.Version)

	if lesson.Status == domain.LessonStatusPublic {
		lesson.SpeechURL = infrastructure.CloudStorageURL + infrastructure.PublicBucketName() + "/" + speechFilePath
		lesson.BodyURL = infrastructure.CloudStorageURL + infrastructure.PublicBucketName() + "/" + bodyFilePath
	} else if lesson.Status == domain.LessonStatusLimited {
		fileType := "" // this is unnecessary when GET request
		bucketName := infrastructure.MaterialBucketName()

		var speechURL string
		var bodyURL string
		var err error

		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			speechURL, err = infrastructure.GetGCSSignedURL(ctx, bucketName, speechFilePath, "GET", fileType)
			return err
		})

		g.Go(func() error {
			bodyURL, err = infrastructure.GetGCSSignedURL(ctx, bucketName, bodyFilePath, "GET", fileType)
			return err
		})

		if err := g.Wait(); err != nil {
			return err
		}

		lesson.SpeechURL = speechURL
		lesson.BodyURL = bodyURL
	}

	return nil
}
