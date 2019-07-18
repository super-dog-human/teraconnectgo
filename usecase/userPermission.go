package usecase

import (
	"context"
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
)

func currentUserAccessToLesson(ctx context.Context, request *http.Request, lessonID string) error {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLessonById(ctx, lessonID, lessonID)
	if err != nil {
		return err
	}

	if lesson.UserID != currentUser.ID {
		return LessonNotAvailable
	}

	return nil
}