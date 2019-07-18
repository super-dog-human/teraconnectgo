package usecase

import (
	"context"
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)

// GetRawVoiceTexts for fetch voice textsfrom Cloud Datastore
func GetRawVoiceTexts(request *http.Request, lessonID string) ([]domain.RawVoiceText, error) {
	ctx := appengine.NewContext(request)

	if err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return nil, err
	}

	voiceTexts, err := domain.GetRawVoiceTexts(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	return voiceTexts, nil
}
