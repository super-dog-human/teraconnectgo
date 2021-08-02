package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateLessonThumbnailBlankFile is create blank image file to public or private bucket.
func CreateLessonThumbnailBlankFile(request *http.Request, isPublic bool, id int64) (string, error) {
	ctx := request.Context()
	url, err := domain.CreateLessonThumbnailBlankFile(ctx, id, isPublic)
	if err != nil {
		return "", err
	}

	return url, nil
}
