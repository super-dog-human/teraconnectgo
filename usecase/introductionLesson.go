package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateIntroductionLesson is create the new lesson and graphics for self-introduction.
func CreateIntroductionLesson(request *http.Request) error {
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return InvalidLessonParams
	}

	ctx := request.Context()
	lessonID, err := domain.CreateIntroductionLesson(ctx, &currentUser)
	if err != nil {
		return err
	}

	if err != domain.CreateIntroductionGraphics(ctx, currentUser.ID, lessonID) {
		return err
	}

	return nil
}
