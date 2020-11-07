package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetRawVoiceTexts for fetch voice textsfrom Cloud Datastore
func GetRawVoiceTexts(request *http.Request, lessonID int64) ([]domain.RawVoiceText, error) {
	ctx := request.Context()
	if err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return nil, err
	}

	voiceTexts, err := domain.GetRawVoiceTexts(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	return voiceTexts, nil
}
