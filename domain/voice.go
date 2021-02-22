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
	URL         string    `json:"url" datastore:"-"`
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

	for i, voice := range *voices {
		voice.ID = keys[i].ID
	}

	if len(*voices) == 0 {
		return VoiceNotFound
	}

	return nil
}

// CreateVoice is creats new voice.
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
