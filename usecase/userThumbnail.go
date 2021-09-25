package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateLessonThumbnailBlankFile is create blank image file.
func CreateUserThumbnailBlankFile(request *http.Request) (string, error) {
	user, err := domain.GetCurrentUser(request)
	if err != nil {
		return "", UserNotAvailable
	}

	ctx := request.Context()
	url, err := domain.CreateUserThumbnailBlankFile(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return url, nil
}
