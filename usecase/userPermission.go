package usecase

import (
	"context"
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

func currentUserAccessToLesson(ctx context.Context, request *http.Request, lessonID int64) error {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lesson, err := domain.GetLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	if lesson.UserID != currentUser.ID {
		return LessonNotAvailable
	}

	return nil
}
