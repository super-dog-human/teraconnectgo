package usecase

import (
	"errors"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/domain"
)

// LessonMaterialParams
type LessonMaterialParams struct {
	DurationSec          float32                     `json:"durationSec"`
	AvatarID             int64                       `json:"avatarID"`
	AvatarLightColor     string                      `json:"avatarLightColor"`
	BackgroundImageID    int64                       `json:"backgroundImageID"`
	VoiceSynthesisConfig domain.VoiceSynthesisConfig `json:"voiceSynthesisConfig"`
	Avatars              []domain.LessonAvatar       `json:"avatars"`
	Drawings             []domain.LessonDrawing      `json:"drawings"`
	Embeddings           []domain.LessonEmbedding    `json:"embeddings"`
	Graphics             []domain.LessonGraphic      `json:"graphics"`
	Musics               []domain.LessonMusic        `json:"musics"`
	Speeches             []domain.LessonSpeech       `json:"speeches"`
}

type LessonMaterialErrorCode uint

const (
	LessonMaterialNotAvailable LessonMaterialErrorCode = 1
	LessonMaterialNotFound     LessonMaterialErrorCode = 2
)

func (e LessonMaterialErrorCode) Error() string {
	switch e {
	case LessonMaterialNotAvailable:
		return "lesson material not available"
	case LessonMaterialNotFound:
		return "lesson material not found"
	default:
		return "unknown lesson error"
	}
}

func GetLessonMaterial(request *http.Request, id int64, lessonID int64) (domain.LessonMaterial, error) {
	ctx := request.Context()

	var lessonMaterial domain.LessonMaterial
	userID, err := currentUserAccessToLesson(ctx, request, lessonID)
	if err != nil {
		return lessonMaterial, LessonMaterialNotAvailable
	}

	if err := domain.GetLessonMaterial(ctx, id, lessonID, &lessonMaterial); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return lessonMaterial, LessonMaterialNotFound
		} else {
			return lessonMaterial, LessonMaterialNotAvailable
		}
	}

	if lessonMaterial.AvatarID != 0 {
		avatar, err := domain.GetPublicAvatarByID(ctx, lessonMaterial.AvatarID)
		if err != nil {
			if ok := errors.Is(err, domain.AvatarNotFound); ok {
				avatar, err = domain.GetCurrentUsersAvatarByID(ctx, lessonMaterial.AvatarID, userID)
				if err != nil {
					return lessonMaterial, err
				}
			} else {
				return lessonMaterial, err
			}
		}

		lessonMaterial.Avatar = avatar
	}

	return lessonMaterial, nil
}

func CreateLessonMaterial(request *http.Request, lessonID int64, params LessonMaterialParams) (int64, error) {
	ctx := request.Context()

	userID, err := currentUserAccessToLesson(ctx, request, lessonID)
	if err != nil {
		return 0, LessonMaterialNotAvailable
	}

	var lessonMaterial domain.LessonMaterial
	copier.Copy(&lessonMaterial, &params)
	lessonMaterial.UserID = userID

	if lessonMaterial.VoiceSynthesisConfig.LanguageCode == "" {
		lessonMaterial.VoiceSynthesisConfig.LanguageCode = "ja-JP"
	}

	if lessonMaterial.VoiceSynthesisConfig.Name == "" {
		lessonMaterial.VoiceSynthesisConfig.Name = "ja-JP-Wavenet-A"
	}

	if err := domain.CreateLessonMaterial(ctx, lessonID, &lessonMaterial); err != nil {
		return 0, err
	}

	lesson, err := domain.GetLessonByID(ctx, lessonID)
	if err != nil {
		return 0, err
	}

	lesson.MaterialID = lessonMaterial.ID

	if err = domain.UpdateLesson(ctx, &lesson); err != nil {
		return 0, err
	}

	return lessonMaterial.ID, nil
}

func UpdateLessonMaterial(request *http.Request, id int64, lessonID int64, params LessonMaterialParams) error {
	ctx := request.Context()

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return LessonMaterialNotAvailable
	}

	var lessonMaterial domain.LessonMaterial
	copier.Copy(&lessonMaterial, &params)

	if err := domain.UpdateLessonMaterial(ctx, id, lessonID, &lessonMaterial); err != nil {
		return err
	}

	return nil
}
