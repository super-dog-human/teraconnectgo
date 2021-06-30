package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateLessonThumbnailBlankFile is create blank image file to public or private bucket.
func CreateLessonThumbnailBlankFile(request *http.Request, id int64) (string, error) {
	ctx := request.Context()

	isPublic := request.URL.Query().Get("is_public") == "true"
	url, err := domain.CreateLessonThumbnailBlankFile(ctx, id, isPublic)
	if err != nil {
		return "", err
	}

	return url, nil
}
