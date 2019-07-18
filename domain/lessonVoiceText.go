package domain

import (
	"context"

	"google.golang.org/appengine/datastore"
)

func GetRawVoiceTexts(ctx context.Context, id string) ([]RawVoiceText, error) {
	var voiceTexts []RawVoiceText
	query := datastore.NewQuery("RawVoiceText").Filter("LessonID =", id)
	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return nil, err
	}

	return voiceTexts, nil
}
