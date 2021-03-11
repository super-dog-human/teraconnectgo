package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type VoiceErrorCode uint

const (
	VoiceNotFound VoiceErrorCode = 1
)

func (e VoiceErrorCode) Error() string {
	switch e {
	case VoiceNotFound:
		return "voice not found"
	default:
		return "unknown voice error"
	}
}

// Voice is used for lesson.
type Voice struct {
	ID          int64     `json:"id" datastore:"-"`
	UserID      int64     `json:"userID"`
	Elapsedtime float32   `json:"elapsedtime"`
	DurationSec float32   `json:"durationSec"`
	Text        string    `json:"text"`
	IsTexted    bool      `json:"isTexted"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// GetVoices is get voice entities belongs to lesson.
func GetVoices(ctx context.Context, lessonID int64, voices *[]Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	query := datastore.NewQuery("Voice").Ancestor(ancestor).Order("Elapsedtime")

	keys, err := client.GetAll(ctx, query, voices)
	if err != nil {
		return err
	}

	if len(*voices) == 0 {
		return VoiceNotFound
	}

	storeVoiceURLs(ctx, lessonID, voices, keys)

	return nil
}

// CreateVoice is creates new voice.
func CreateVoice(ctx context.Context, lessonID int64, voice *Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	voice.Created = time.Now()

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IncompleteKey("Voice", ancestor)
	putKey, err := client.Put(ctx, key, voice)
	if err != nil {
		return err
	}

	voice.ID = putKey.ID

	return nil
}

func storeVoiceURLs(ctx context.Context, lessonID int64, voices *[]Voice, keys []*datastore.Key) error {
	//	lessonIDString := strconv.FormatInt(lessonID, 10)
	for i, key := range keys {
		/*
			fileID := strconv.FormatInt(key.ID, 10)
			filePath := fmt.Sprintf("voice/%s/%s.mp3", lessonIDString, fileID)
			fileType := "" // this is unnecessary when GET request
			bucketName := infrastructure.MaterialBucketName()

			url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)
			if err != nil {
				return err
			}
		*/

		(*voices)[i].ID = key.ID
		//		(*voices)[i].URL = url
	}

	return nil
}
