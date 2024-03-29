package usecase

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/domain"
)

// NewLessonMaterialParams
type NewLessonMaterialParams struct {
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

func UpdateLessonMaterial(request *http.Request, id int64, lessonID int64, params *map[string]interface{}) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLessonByID(ctx, lessonID)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return LessonNotFound
		}
		return err
	}

	if lesson.UserID != currentUser.ID {
		return LessonMaterialNotAvailable
	}

	if id != lesson.MaterialID {
		return LessonMaterialNotAvailable
	}

	targetFields := []string{"DurationSec", "Avatars", "Drawings", "Embeddings", "Graphics", "Musics", "Speeches"}
	if err := domain.UpdateLessonMaterial(ctx, id, lessonID, params, &targetFields); err != nil {
		return err
	}

	return nil
}

func createInitialLessonMaterial(ctx context.Context, userID int64, lessonID int64) (int64, error) {
	var materialID int64

	avatars, err := domain.GetPublicAvatars(ctx) // 数が少ないので全件取得して1件使用する
	if err != nil {
		return materialID, err
	}
	avatarID := avatars[0].ID

	backgroundImage, err := domain.GetBackgroundImage(ctx)
	if err != nil {
		return materialID, err
	}

	materialID, err = domain.CreateInitialLessonMaterial(ctx, userID, avatarID, backgroundImage.ID, lessonID)
	if err != nil {
		return materialID, err
	}

	return materialID, nil
}
