package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

func GetLessonMaterial(request *http.Request, lessonID int64) (domain.LessonMaterial, error) {
	ctx := request.Context()

	var lessonMaterial domain.LessonMaterial
	if err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return lessonMaterial, err
	}

	lessonMaterial, err := domain.GetLessonMaterialFromGCS(ctx, lessonID)
	if err != nil {
		return lessonMaterial, err
	}

	return lessonMaterial, nil
}

func CreateLessonMaterial(request *http.Request, lessonID int64, lessonMaterial domain.LessonMaterial) error {
	ctx := request.Context()

	if err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return err
	}

	if err := domain.CreateLessonMaterialFileToGCS(ctx, lessonID, lessonMaterial); err != nil {
		return err
	}

	return nil
}

func UpdateLessonMaterial(request *http.Request, lessonID int64, lessonMaterial domain.LessonMaterial) error {
	ctx := request.Context()

	if err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return err
	}

	if err := domain.CreateLessonMaterialFileToGCS(ctx, lessonID, lessonMaterial); err != nil {
		return err
	}

	return nil
}
