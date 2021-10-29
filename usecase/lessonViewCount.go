package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// UpdateLessonViewCountは、lessonIDのLessonのViewCountを増分します。
func UpdateLessonViewCount(request *http.Request, lessonID int64) error {
	ctx := request.Context()

	if infrastructure.AppEnv() == "production" {
		if err := domain.IncrementLessonViewCount(ctx, lessonID); err != nil {
			return err
		}
	}

	return nil
}
