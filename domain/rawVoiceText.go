package domain

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetRawVoiceTexts(ctx context.Context, lessonID int64) ([]RawVoiceText, error) {
	var voiceTexts []RawVoiceText

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("RawVoiceText").Filter("LessonID =", lessonID).Order("FileID")
	if _, err := client.GetAll(ctx, query, &voiceTexts); err != nil {
		return voiceTexts, err
	}

	return voiceTexts, nil
}

func DeleteRawVoiceTextsByLessonID(ctx context.Context, lessonID int64) error {
	var voiceTexts []RawVoiceText
	var keys []*datastore.Key

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	query := datastore.NewQuery("RawVoiceText").Filter("LessonID =", lessonID)
	if _, err := client.GetAll(ctx, query, &voiceTexts); err != nil {
		return err
	}

	if err := client.DeleteMulti(ctx, keys); err != nil {
		return err
	}

	return nil
}

func DeleteRawVoiceTextsInTransactionByLessonID(tx *datastore.Transaction, voiceTexts []RawVoiceText) error {
	var keys []*datastore.Key

	for _, voiceText := range voiceTexts {
		keys = append(keys, datastore.NameKey("RawVoiceText", voiceText.FileID, nil))
	}

	if err := tx.DeleteMulti(keys); err != nil {
		return err
	}

	return nil
}
