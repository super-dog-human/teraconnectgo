package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// PackLesson is packing the lesson to zip and upload GCS
func PackLesson(request *http.Request, id int64) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLessonByID(ctx, id)
	if err != nil {
		return err
	}

	if lesson.UserID != currentUser.ID {
		return LessonNotAvailable
	}

	/*

		graphicFileTypes, err := domain.GetGraphicFileTypes(ctx, lesson.GraphicIDs)
		if err != nil {
			return err
		}

		voiceTexts, err := domain.GetRawVoiceTexts(ctx, id)
		if err != nil {
			return err
		}

		zip, err := domain.CreateLessonZip(ctx, lesson, graphicFileTypes, voiceTexts)
		if err != nil {
			return err
		}

		if err = domain.UploadLessonZipToGCS(ctx, id, zip); err != nil {
			return err
		}

		lesson.IsPacked = true
		if err = domain.UpdateLesson(ctx, &lesson); err != nil {
			return err
		}

		if err = domain.CreateLessonIndex(ctx, currentUser, &lesson, &voiceTexts); err != nil {
			return err
		}
	*/

	return nil
}
