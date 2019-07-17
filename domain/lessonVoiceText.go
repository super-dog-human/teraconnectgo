package domain

import (
	"context"

	"google.golang.org/appengine/datastore"
)

func GetLessonVoiceTexts(ctx context.Context, id string) ([]LessonVoiceText, error) {
	var voiceTexts []LessonVoiceText
	query := datastore.NewQuery("LessonVoiceText").Filter("LessonID =", id)
	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return nil, err
	}

	return voiceTexts, nil
}
