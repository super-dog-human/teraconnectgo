package usecase

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)

// PackLesson is packing the lesson to zip and upload GCS
func PackLesson(request *http.Request, id string) error {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLesson(request, id)
	if err != nil {
		return err
	}

	if lesson.UserID != currentUser.ID {
		return LessonNotAvailable
	}

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

	if err = domain.RemoveUsedFilesInGCS(ctx, id, voiceTexts); err != nil {
		return err
	}

	lesson.IsPacked = true
	if err = domain.UpdateLesson(ctx, lesson); err != nil {
		return err
	}

	return nil
}
