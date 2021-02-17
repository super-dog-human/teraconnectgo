package usecase

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateLessonMaterialParams
type CreateLessonMaterialParams struct {
	Duration          float32                     `json:"duration"`
	AvatarID          int64                       `json:"avatarID"`
	AvatarLightColor  string                      `json:"avatarLightColor"`
	BackgroundImageID int64                       `json:"backgroundImageID"`
	BackgroundMusicID int64                       `json:"backgroundMusicID"`
	AvatarMovings     []domain.LessonAvatarMoving `json:"avatarMovings"`
	Graphics          []domain.LessonGraphic      `json:"graphics"`
	Drawings          []domain.LessonDrawing      `json:"drawings"`
}

func GetLessonMaterial(request *http.Request, lessonID int64) (domain.LessonMaterial, error) {
	ctx := request.Context()

	var lessonMaterial domain.LessonMaterial
	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return lessonMaterial, err
	}
	/*
		lessonMaterial, err := domain.GetLessonMaterialFromGCS(ctx, lessonID)
		if err != nil {
			return lessonMaterial, err
		}

	*/
	return lessonMaterial, nil
}

func CreateLessonMaterial(request *http.Request, lessonID int64, params CreateLessonMaterialParams) error {
	ctx := request.Context()

	userID, err := currentUserAccessToLesson(ctx, request, lessonID)
	if err != nil {
		return err
	}

	var lessonMaterial domain.LessonMaterial
	copier.Copy(&lessonMaterial, &params)
	lessonMaterial.UserID = userID

	if err := domain.CreateLessonMaterial(ctx, lessonID, &lessonMaterial); err != nil {
		return err
	}

	return nil
}
